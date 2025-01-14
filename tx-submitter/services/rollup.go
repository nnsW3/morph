package services

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/holiman/uint256"
	"github.com/scroll-tech/go-ethereum"
	"github.com/scroll-tech/go-ethereum/accounts/abi"
	"github.com/scroll-tech/go-ethereum/accounts/abi/bind"
	"github.com/scroll-tech/go-ethereum/common"
	"github.com/scroll-tech/go-ethereum/consensus/misc/eip4844"
	"github.com/scroll-tech/go-ethereum/core"
	"github.com/scroll-tech/go-ethereum/core/types"
	"github.com/scroll-tech/go-ethereum/crypto"
	"github.com/scroll-tech/go-ethereum/crypto/bls12381"
	"github.com/scroll-tech/go-ethereum/eth"
	"github.com/scroll-tech/go-ethereum/log"
	"github.com/scroll-tech/go-ethereum/params"
	"github.com/scroll-tech/go-ethereum/rpc"
	"github.com/tendermint/tendermint/blssignatures"

	"morph-l2/bindings/bindings"
	"morph-l2/tx-submitter/iface"
	"morph-l2/tx-submitter/localpool"
	"morph-l2/tx-submitter/metrics"
	"morph-l2/tx-submitter/utils"
)

const (
	txSlotSize  = 32 * 1024
	txMaxSize   = 4 * txSlotSize // 128KB
	rotatorWait = 3 * time.Second
	rotatorBuff = 15
)

type Rollup struct {
	ctx     context.Context
	metrics *metrics.Metrics

	l1RpcClient *rpc.Client
	L1Client    iface.Client
	L2Clients   []iface.L2Client
	Rollup      iface.IRollup

	Staking iface.IL1Staking

	chainId    *big.Int
	privKey    *ecdsa.PrivateKey
	rollupAddr common.Address
	abi        *abi.ABI

	// rotator
	rotator    *Rotator
	pendingTxs *PendingTxs

	rollupFinalizeMu sync.Mutex

	// cfg
	cfg utils.Config
}

func NewRollup(
	ctx context.Context,
	metrics *metrics.Metrics,
	l1RpcClient *rpc.Client,
	l1 iface.Client,
	l2Clients []iface.L2Client,
	rollup iface.IRollup,
	staking iface.IL1Staking,
	chainId *big.Int,
	priKey *ecdsa.PrivateKey,
	rollupAddr common.Address,
	abi *abi.ABI,
	cfg utils.Config,
	rotator *Rotator,
) *Rollup {

	return &Rollup{
		ctx:         ctx,
		metrics:     metrics,
		l1RpcClient: l1RpcClient,
		L1Client:    l1,
		Rollup:      rollup,
		Staking:     staking,
		L2Clients:   l2Clients,
		privKey:     priKey,
		chainId:     chainId,
		rollupAddr:  rollupAddr,
		abi:         abi,
		rotator:     rotator,
		cfg:         cfg,
	}
}

func (sr *Rollup) Start() {

	// journal
	jn := localpool.New(sr.cfg.JournalFilePath)
	err := jn.Init()
	if err != nil {
		log.Crit("journal file init failed", "err", err)
	}
	// pendingtxs
	sr.pendingTxs = NewPendingTxs(sr.abi.Methods["commitBatch"].ID, sr.abi.Methods["finalizeBatch"].ID, jn)
	txs, err := jn.ParseAllTxs()
	if err != nil {
		log.Error("parse l1 mempool error", "error", err)
	} else {
		sr.pendingTxs.Recover(txs, sr.abi)
	}

	// metrics
	go utils.Loop(sr.ctx, 10*time.Second, func() {

		// get balacnce of wallet
		balance, err := sr.L1Client.BalanceAt(context.Background(), crypto.PubkeyToAddress(sr.privKey.PublicKey), nil)
		if err != nil {
			log.Error("get wallet balance error", "error", err)
			if utils.IsRpcErr(err) {
				sr.metrics.IncRpcErrors()
			}
			return
		}
		// balance to eth
		balanceEth := new(big.Rat).SetFrac(balance, big.NewInt(params.Ether))

		// parse float64 from string
		balanceEthFloat, err := strconv.ParseFloat(balanceEth.FloatString(18), 64)
		if err != nil {
			log.Warn("parse balance to float error", "error", err)
			return
		}

		sr.metrics.SetWalletBalance(balanceEthFloat)

	})

	go utils.Loop(sr.ctx, sr.cfg.RollupInterval, func() {
		sr.rollupFinalizeMu.Lock()
		defer sr.rollupFinalizeMu.Unlock()
		if err := sr.rollup(); err != nil {
			if utils.IsRpcErr(err) {
				sr.metrics.IncRpcErrors()
			}
			log.Error("rollup failed,wait for the next try", "error", err)
		}
	})

	if sr.cfg.Finalize {

		go utils.Loop(sr.ctx, sr.cfg.FinalizeInterval, func() {
			sr.rollupFinalizeMu.Lock()
			defer sr.rollupFinalizeMu.Unlock()

			if err := sr.finalize(); err != nil {
				log.Error("finalize failed", "error", err)
				if utils.IsRpcErr(err) {
					sr.metrics.IncRpcErrors()
				}

			}
		})
	}

	var processtxMu sync.Mutex
	go utils.Loop(sr.ctx, sr.cfg.TxProcessInterval, func() {
		processtxMu.Lock()
		defer processtxMu.Unlock()
		if err := sr.ProcessTx(); err != nil {
			log.Error("process tx err", "error", err)
			if utils.IsRpcErr(err) {
				sr.metrics.IncRpcErrors()
			}
		}
	})

}

func (sr *Rollup) ProcessTx() error {

	// case 1: in mempool
	//          -> check timeout
	// case 2: no in mempool
	// case 2.1: discarded
	// case 2.2: tx included -> success
	// case 2.3: tx included -> failed
	//          -> reset index to failed index

	// get all local txs
	txRecords := sr.pendingTxs.GetAll()
	if len(txRecords) == 0 {
		return nil
	}

	// query tx status
	for _, txRecord := range txRecords {

		rtx := txRecord.tx
		method := utils.ParseMethod(rtx, sr.abi)
		log.Info("process tx", "hash", rtx.Hash().String(), "nonce", rtx.Nonce(), "method", method)
		// query tx
		_, ispending, err := sr.L1Client.TransactionByHash(context.Background(), txRecord.tx.Hash())
		if err != nil {
			if !utils.ErrStringMatch(err, ethereum.NotFound) {
				return fmt.Errorf("query tx  error:%w, tx: %s, nonce: %d", err, rtx.Hash().String(), rtx.Nonce())
			}
			sr.pendingTxs.IncQueryTimes(rtx.Hash()) // not found in mempool, increase query times
		} else {
			log.Info("query tx success", "hash", rtx.Hash().Hex(), "pending", ispending)
		}

		// exist in mempool
		if ispending {
			if txRecord.sendTime+uint64(sr.cfg.TxTimeout.Seconds()) < uint64(time.Now().Unix()) {
				log.Info("tx timeout", "tx", rtx.Hash().Hex(), "nonce", rtx.Nonce(), "method", method)
				newtx, err := sr.ReSubmitTx(false, &rtx)
				if err != nil {
					log.Error("resubmit tx", "error", err, "tx", rtx.Hash().Hex(), "nonce", rtx.Nonce())
					return fmt.Errorf("resubmit tx error:%w", err)
				} else {
					log.Info("replace success", "old_tx", rtx.Hash().Hex(), "new_tx", newtx.Hash().String(), "nonce", rtx.Nonce())
					sr.pendingTxs.Remove(rtx.Hash())
					sr.pendingTxs.Add(*newtx)
				}
			}
		} else { // not in mempool
			receipt, err := sr.L1Client.TransactionReceipt(context.Background(), rtx.Hash())
			if err != nil {
				log.Error("query tx receipt error", "tx", rtx.Hash().String(), "nonce", rtx.Nonce(), "error", err)
				if !utils.ErrStringMatch(err, ethereum.NotFound) {
					return err
				}

				// sr.pendingTxs.txinfos
				if txRecord.queryTimes >= 5 {
					log.Warn("tx discarded",
						"hash", rtx.Hash().String(),
						"nonce", rtx.Nonce(),
						"query_times", txRecord.queryTimes,
					)
					replacedtx, err := sr.ReSubmitTx(true, &rtx)
					if err != nil {
						log.Error("resend discarded tx", "old_tx", rtx.Hash().String(), "nonce", rtx.Nonce(), "error", err)
						if utils.ErrStringMatch(err, core.ErrNonceTooLow) {
							log.Info("discarded tx removed",
								"hash", rtx.Hash().String(),
								"nonce", rtx.Nonce(),
								"method", method,
							)
							sr.pendingTxs.Remove(rtx.Hash())
							return nil
						}
						return fmt.Errorf("resend discarded tx: %w", err)
					} else {
						sr.pendingTxs.Remove(rtx.Hash())
					}
					sr.pendingTxs.Add(*replacedtx)
					log.Info("resend discarded tx", "old_tx", rtx.Hash().String(), "new_tx", replacedtx.Hash().String(), "nonce", replacedtx.Nonce())
				} else {
					log.Info("tx is not found, neither in mempool nor in block", "hash", rtx.Hash().String(), "nonce", rtx.Nonce(), "query_times", txRecord.queryTimes)
				}
			} else {
				logs := utils.ParseBusinessInfo(rtx, sr.abi)
				logs = append(logs,
					"block", receipt.BlockNumber,
					"hash", rtx.Hash().String(),
					"status", receipt.Status,
					"gas_used", receipt.GasUsed,
					"type", rtx.Type(),
					"nonce", rtx.Nonce(),
					"blob_fee_cap", rtx.BlobGasFeeCap(),
					"blob_gas", rtx.BlobGas(),
					"tx_size", rtx.Size(),
					"gas_limit", rtx.Gas(),
					"gas_price", rtx.GasPrice(),
				)

				log.Info("tx included",
					logs...,
				)

				if receipt.Status != types.ReceiptStatusSuccessful {
					// if tx is commitBatch
					if method == "commitBatch" {
						parentindex := utils.ParseParentBatchIndex(rtx.Data())
						index := parentindex + 1

						// prevent the SetFailedStatus operation from
						// happening between RemoveRollupRestriction
						// and SetPindex in the rollup function
						sr.rollupFinalizeMu.Lock()
						sr.pendingTxs.SetFailedStatus(index)
						sr.rollupFinalizeMu.Unlock()

					}

				} else {
					if method == "commitBatch" && sr.pendingTxs.failedIndex != nil {
						log.Info("fail revover", "failed_index", sr.pendingTxs.failedIndex)
						sr.pendingTxs.RemoveRollupRestriction()
					}
				}

				sr.pendingTxs.Remove(rtx.Hash())
				// set metrics
				fee := calcFee(receipt)
				if fee == 0 {
					log.Warn("fee is zero", "hash", rtx.Hash().Hex())
				}
				if method == "commitBatch" {
					sr.metrics.SetRollupCost(fee)
				} else if method == "finalizeBatch" {
					sr.metrics.SetFinalizeCost(fee)
				}
			}

		}

	}

	return nil

}

func (sr *Rollup) finalize() error {
	// get last finalized
	lastFinalized, err := sr.Rollup.LastFinalizedBatchIndex(nil)
	if err != nil {
		return fmt.Errorf("get last finalized error:%v", err)
	}
	// get last committed
	lastCommitted, err := sr.Rollup.LastCommittedBatchIndex(nil)
	if err != nil {
		return fmt.Errorf("get last committed error:%v", err)
	}

	target := big.NewInt(int64(sr.pendingTxs.pfinalize + 1))
	if target.Cmp(lastFinalized) <= 0 {
		target = new(big.Int).Add(lastFinalized, big.NewInt(1))
	}

	if target.Cmp(lastCommitted) > 0 {
		log.Info("no need to finalize", "last_finalized", lastFinalized.Uint64(), "last_committed", lastCommitted.Uint64())
		return nil
	}

	log.Info("finalize info",
		"last_fianlzied", lastFinalized,
		"last_committed", lastCommitted,
		"finalize_index", target,
	)

	// batch exist
	existed, err := sr.Rollup.BatchExist(nil, target)
	if err != nil {
		log.Error("query batch exist", "err", err)
		return err
	}
	if !existed {
		log.Warn("finalized batch not existed")
		return nil
	}

	// in challenge window
	inWindow, err := sr.Rollup.BatchInsideChallengeWindow(nil, target)
	if err != nil {
		return fmt.Errorf("get batch inside challenge window error:%v", err)
	}
	if inWindow {
		log.Info("batch inside challenge window, wait")
		return nil
	}
	// finalize
	opts, err := bind.NewKeyedTransactorWithChainID(sr.privKey, sr.chainId)
	if err != nil {
		return fmt.Errorf("new keyedTransaction with chain id error:%v", err)
	}

	// get next batch
	nextBatchIndex := target.Uint64() + 1
	batch, err := GetRollupBatchByIndex(nextBatchIndex, sr.L2Clients)
	if err != nil {
		log.Error("get next batch by index error",
			"batch_index", nextBatchIndex,
		)
		return fmt.Errorf("get next batch by index err:%v", err)
	}
	if batch == nil {
		log.Info("next batch is nil,wait next batch header to finalize", "next_batch_index", nextBatchIndex)
		return nil
	}

	// calldata
	calldata, err := sr.abi.Pack("finalizeBatch", []byte(batch.ParentBatchHeader))
	if err != nil {
		return fmt.Errorf("pack finalizeBatch error:%v", err)
	}
	tip, feecap, _, err := sr.GetGasTipAndCap()
	if err != nil {
		log.Error("get gas tip and cap error", "business", "finalize")
		return fmt.Errorf("get gas tip and cap error:%v", err)
	}

	gasDefault := uint64(50_0000)
	gas, err := sr.EstimateGas(opts.From, sr.rollupAddr, calldata, feecap, tip)
	if err != nil {
		gas = gasDefault
	} else {
		gas = gas * 12 / 10 // add a buffer
	}

	var nonce uint64
	if sr.pendingTxs.pnonce != 0 {
		nonce = sr.pendingTxs.pnonce + 1
	} else {
		nonce, err = sr.L1Client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(sr.privKey.PublicKey))
		if err != nil {
			return fmt.Errorf("query layer1 nonce error:%v", err.Error())
		}
	}

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   sr.chainId,
		Nonce:     nonce,
		GasTipCap: tip,
		GasFeeCap: feecap,
		Gas:       gas,
		To:        &sr.rollupAddr,
		Data:      calldata,
	})

	if uint64(tx.Size()) > txMaxSize {
		return core.ErrOversizedData
	}

	signedTx, err := opts.Signer(opts.From, tx)
	if err != nil {
		return fmt.Errorf("sign tx error:%v", err)
	}

	log.Info("finalize tx info",
		"batch_index", target,
		"last_committed", lastCommitted,
		"last_finalized", lastFinalized,
		"hash", signedTx.Hash().String(),
		"type", signedTx.Type(),
		"nonce", signedTx.Nonce(),
		"gas", signedTx.Gas(),
		"tip", signedTx.GasTipCap().String(),
		"fee_cap", signedTx.GasFeeCap().String(),
		"size", signedTx.Size(),
	)

	err = sr.SendTx(signedTx)
	if err != nil {
		log.Error("send finalize tx to mempool", "error", err.Error())
		if utils.ErrStringMatch(err, core.ErrNonceTooLow) {
			// adjust nonce
			n1, _, err := utils.ParseNonce(err.Error())
			if err != nil {
				return fmt.Errorf("parse nonce err: %w", err)
			}
			sr.pendingTxs.SetNonce(n1 - 1)
			log.Info("update pnonce", "nonce", n1-1)
		}
		return fmt.Errorf("send tx error:%v", err.Error())
	} else {
		log.Info("finalzie tx sent")

		sr.pendingTxs.SetNonce(signedTx.Nonce())
		sr.pendingTxs.SetPFinalize(target.Uint64())
		sr.pendingTxs.Add(*signedTx)
	}

	return nil

}

func (sr *Rollup) rollup() error {

	if !sr.cfg.PriorityRollup {
		cur, err := sr.rotator.CurrentSubmitter(sr.L2Clients)
		if err != nil {
			return fmt.Errorf("rollup: get current submitter err, %w", err)
		}

		past := (time.Now().Unix() - sr.rotator.GetStartTime().Int64()) % sr.rotator.GetEpoch().Int64()
		start := time.Now().Unix() - past
		end := start + sr.rotator.GetEpoch().Int64()

		log.Info("rotator info",
			"turn", cur.Hex(),
			"cur", sr.walletAddr(),
			"start", start,
			"end", end,
			"now", time.Now().Unix(),
		)

		if cur.Hex() == sr.walletAddr() {
			left := end - time.Now().Unix()
			if left < rotatorBuff {
				log.Info("rollup time not enough, wait next turn", "left", left)
				return nil
			}

			log.Info("start to rollup")
		} else {
			log.Info("wait for my turn")
			return nil
		}
	}

	if len(sr.pendingTxs.txinfos) > int(sr.cfg.MaxTxsInPendingPool) {
		log.Info("too many txs in mempool, wait")
		return nil
	}

	var nonce uint64
	var batchIndex uint64

	cindexBig, err := sr.Rollup.LastCommittedBatchIndex(nil)
	if err != nil {
		return fmt.Errorf("get last committed batch index error:%v", err)
	}
	cindex := cindexBig.Uint64()

	if sr.pendingTxs.failedIndex != nil && cindex >= *sr.pendingTxs.failedIndex {
		sr.pendingTxs.RemoveRollupRestriction()
	}

	if sr.pendingTxs.failedIndex != nil {
		batchIndex = *sr.pendingTxs.failedIndex
	} else {
		if sr.pendingTxs.pindex != 0 {
			if cindex > sr.pendingTxs.pindex {
				batchIndex = cindex + 1
			} else {
				batchIndex = sr.pendingTxs.pindex + 1
			}

		} else {
			batchIndex = cindex + 1
		}
	}

	log.Info("batch info", "last_commit_batch", batchIndex-1, "batch_will_get", batchIndex)
	if sr.pendingTxs.ExistedIndex(batchIndex) {
		log.Info("batch index already committed", "index", batchIndex)
		return nil
	}

	if sr.pendingTxs.failedIndex != nil && batchIndex > *sr.pendingTxs.failedIndex {
		log.Warn("rollup rejected", "index", batchIndex)
		return nil
	}

	batch, err := GetRollupBatchByIndex(batchIndex, sr.L2Clients)
	if err != nil {
		return fmt.Errorf("get rollup batch by index err:%v", err)
	}

	// check if the batch is valid
	if batch == nil {
		log.Info("new batch not found, wait for the next turn")
		return nil
	}

	if len(batch.Signatures) == 0 {
		log.Info("length of batch signature is empty, wait for the next turn")
		return nil
	}

	var chunks [][]byte
	// var blobChunk []byte
	for _, chunk := range batch.Chunks {
		chunks = append(chunks, chunk)
	}

	signature, err := sr.aggregateSignatures(batch)
	if err != nil {
		return err
	}
	rollupBatch := bindings.IRollupBatchDataInput{
		Version:                uint8(batch.Version),
		ParentBatchHeader:      batch.ParentBatchHeader,
		Chunks:                 chunks,
		SkippedL1MessageBitmap: batch.SkippedL1MessageBitmap,
		PrevStateRoot:          batch.PrevStateRoot,
		PostStateRoot:          batch.PostStateRoot,
		WithdrawalRoot:         batch.WithdrawRoot,
	}

	opts, err := bind.NewKeyedTransactorWithChainID(sr.privKey, sr.chainId)
	if err != nil {
		return fmt.Errorf("new keyedTransaction with chain id error:%v", err)
	}

	// tip and cap
	tip, gasFeeCap, blobFee, err := sr.GetGasTipAndCap()
	if err != nil {
		return fmt.Errorf("get gas tip and cap error:%v", err)
	}

	// calldata encode
	calldata, err := sr.abi.Pack("commitBatch", rollupBatch, *signature)
	if err != nil {
		return fmt.Errorf("pack calldata error:%v", err)
	}

	gas, err := sr.EstimateGas(opts.From, sr.rollupAddr, calldata, gasFeeCap, tip)
	if err != nil {
		log.Warn("estimate gas error", "err", err)
		if sr.pendingTxs.HaveFailed() {
			log.Warn("estimate gas err, wait",
				"err", err,
				"update_index", cindex+1,
			)
			sr.pendingTxs.ResetFailedIndex(cindex + 1)
			return nil
		} else {
			msgcnt := utils.ParseL1MessageCnt(batch.Chunks)
			gas = sr.RoughEstimateGas(msgcnt)
		}
	}

	// gas buffer
	gas = gas * sr.cfg.GasLimitBuffer / 100

	if sr.pendingTxs.pnonce != 0 {
		nonce = sr.pendingTxs.pnonce + 1
	} else {
		nonce, err = sr.L1Client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(sr.privKey.PublicKey))
		if err != nil {
			return fmt.Errorf("query layer1 nonce error:%v", err.Error())
		}
	}

	var tx *types.Transaction
	if len(batch.Sidecar.Blobs) > 0 {
		versionedHashes := make([]common.Hash, 0)
		for _, commit := range batch.Sidecar.Commitments {
			versionedHashes = append(versionedHashes, kZGToVersionedHash(commit))
		}
		// blob tx
		tx = types.NewTx(&types.BlobTx{
			ChainID:    uint256.MustFromBig(sr.chainId),
			Nonce:      nonce,
			GasTipCap:  uint256.MustFromBig(big.NewInt(tip.Int64())),
			GasFeeCap:  uint256.MustFromBig(big.NewInt(gasFeeCap.Int64())),
			Gas:        gas,
			To:         sr.rollupAddr,
			Data:       calldata,
			BlobFeeCap: uint256.MustFromBig(blobFee),
			BlobHashes: versionedHashes,
			Sidecar: &types.BlobTxSidecar{
				Blobs:       batch.Sidecar.Blobs,
				Commitments: batch.Sidecar.Commitments,
				Proofs:      batch.Sidecar.Proofs,
			},
		})

	} else {
		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   sr.chainId,
			Nonce:     nonce,
			GasTipCap: tip,
			GasFeeCap: gasFeeCap,
			Gas:       gas,
			To:        &sr.rollupAddr,
			Data:      calldata,
		})
	}

	opts.Nonce = big.NewInt(int64(nonce))
	var signedTx *types.Transaction
	if tx.Type() == types.BlobTxType {
		signedTx, err = types.SignTx(tx, types.NewLondonSignerWithEIP4844(sr.chainId), sr.privKey)
		if err != nil {
			return fmt.Errorf("sign tx error:%v", err)
		}
	} else {
		signedTx, err = opts.Signer(opts.From, tx)
		if err != nil {
			return fmt.Errorf("sign tx error:%v", err)
		}
	}

	log.Info("rollup tx info",
		"batch_index", batchIndex,
		"hash", signedTx.Hash().String(),
		"type", signedTx.Type(),
		"nonce", signedTx.Nonce(),
		"gas", signedTx.Gas(),
		"tip", signedTx.GasTipCap().String(),
		"fee_cap", signedTx.GasFeeCap().String(),
		"blob_fee_cap", signedTx.BlobGasFeeCap(),
		"blob_gas", signedTx.BlobGas(),
		"size", signedTx.Size(),
		"blob_len", len(signedTx.BlobHashes()),
	)

	err = sr.SendTx(signedTx)
	if err != nil {
		log.Error("send tx to mempool", "error", err.Error())
		if utils.ErrStringMatch(err, core.ErrNonceTooLow) {
			// adjust nonce
			n1, _, err := utils.ParseNonce(err.Error())
			if err != nil {
				return fmt.Errorf("parse nonce err: %w", err)
			}
			sr.pendingTxs.SetNonce(n1 - 1)
			log.Info("update pnonce", "nonce", n1-1)
		}
		return fmt.Errorf("send tx error:%v", err.Error())
	} else {
		log.Info("rollup tx send to mempool succuess", "hash", signedTx.Hash().String())

		sr.pendingTxs.SetPindex(batchIndex)
		sr.pendingTxs.SetNonce(tx.Nonce())
		sr.pendingTxs.Add(*signedTx)
	}

	return nil
}

func (sr *Rollup) aggregateSignatures(batch *eth.RPCRollupBatch) (*bindings.IRollupBatchSignatureInput, error) {
	blsSignatures := batch.Signatures
	if len(blsSignatures) == 0 {
		return nil, fmt.Errorf("invalid batch signature")
	}
	signers := make([]common.Address, len(blsSignatures))
	sigs := make([]blssignatures.Signature, 0)
	for i, bz := range blsSignatures {
		if len(bz.Signature) > 0 {
			sig, err := blssignatures.SignatureFromBytes(bz.Signature)
			if err != nil {
				return nil, err
			}
			sigs = append(sigs, sig)
			signers[i] = bz.Signer
		}
	}
	aggregatedSig := blssignatures.AggregateSignatures(sigs)
	blsSignature := bls12381.NewG1().EncodePoint(aggregatedSig)
	// abi pack
	AddressArr, _ := abi.NewType("address[]", "", nil)
	args := abi.Arguments{
		{Type: AddressArr, Name: "stakeAddresses"},
	}
	bsSigner, err := args.Pack(&signers)
	if err != nil {
		return nil, fmt.Errorf("pack signers error:%v", err)
	}

	sigData := bindings.IRollupBatchSignatureInput{
		SignedSequencers: bsSigner,
		SequencerSets:    batch.CurrentSequencerSetBytes,
		Signature:        blsSignature,
	}
	return &sigData, nil
}

func (sr *Rollup) GetGasTipAndCap() (*big.Int, *big.Int, *big.Int, error) {
	tip, err := sr.L1Client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}
	head, err := sr.L1Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, nil, nil, err
	}
	var gasFeeCap *big.Int
	if head.BaseFee != nil {
		gasFeeCap = new(big.Int).Add(
			tip,
			new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
		)
	} else {
		gasFeeCap = new(big.Int).Set(tip)
	}

	// calc blob fee cap
	var blobFee *big.Int
	if head.ExcessBlobGas != nil {
		blobFee = eip4844.CalcBlobFee(*head.ExcessBlobGas)
	}

	//calldata fee bump x*fee/100
	if sr.cfg.CalldataFeeBump > 0 {
		// feecap
		gasFeeCap = new(big.Int).Mul(gasFeeCap, big.NewInt(int64(sr.cfg.CalldataFeeBump)))
		gasFeeCap = new(big.Int).Div(gasFeeCap, big.NewInt(100))
		// tip
		tip = new(big.Int).Mul(tip, big.NewInt(int64(sr.cfg.CalldataFeeBump)))
		tip = new(big.Int).Div(tip, big.NewInt(100))
	}

	return tip, gasFeeCap, blobFee, nil
}

// Init is run before the submitter to check whether the submitter can be started
func (sr *Rollup) Init() error {

	isStaker, err := sr.IsStaker()
	if err != nil {
		return fmt.Errorf("check if this account is sequencer error:%v", err)
	}

	if !isStaker {
		return fmt.Errorf("this account is not staker, can not rollup")
	}

	return nil
}

func (sr *Rollup) walletAddr() string {
	return crypto.PubkeyToAddress(sr.privKey.PublicKey).Hex()
}

func GetRollupBatchByIndex(index uint64, clients []iface.L2Client) (*eth.RPCRollupBatch, error) {
	if len(clients) < 1 {
		return nil, fmt.Errorf("no client to query")
	}
	for _, client := range clients {
		batch, err := client.GetRollupBatchByIndex(context.Background(), index)
		if err != nil {
			log.Warn("failed to get batch", "error", err)
			continue
		}
		if batch != nil && len(batch.Signatures) > 0 {
			return batch, nil
		}
	}

	return nil, nil
}

// query sequencer set from sequencer contract on l2
func GetSequencerSet(addr common.Address, clients []iface.L2Client) ([]common.Address, error) {
	if len(clients) < 1 {
		return nil, fmt.Errorf("no client to query sequencer set")
	}
	for _, client := range clients {
		// l2 sequencer
		l2Seqencer, err := bindings.NewSequencer(addr, client)
		if err != nil {
			log.Warn("failed to connect to sequencer", "error", err)
			continue
		}
		// get sequencer set
		seqSet, err := l2Seqencer.GetSequencerSet2(nil)
		if err != nil {
			log.Warn("failed to get sequencer set", "error", err)
			continue
		}
		return seqSet, nil
	}
	return nil, fmt.Errorf("no sequencer set found after querying all clients")
}

// query epoch from gov contract on l2
func GetEpoch(addr common.Address, clients []iface.L2Client) (*big.Int, error) {
	if len(clients) < 1 {
		return nil, fmt.Errorf("no client to query epoch")
	}
	for _, client := range clients {
		// l2 gov
		l2Gov, err := bindings.NewGov(addr, client)
		if err != nil {
			log.Warn("failed to connect to gov", "error", err)
			continue
		}
		// get epoch
		epoch, err := l2Gov.RollupEpoch(nil)
		if err != nil {
			log.Warn("failed to get epoch", "error", err)
			continue
		}
		return epoch, nil
	}
	return nil, fmt.Errorf("no epoch found after querying all clients")
}

// query sequencer set update time from sequencer contract on l2
func GetSequencerSetUpdateTime(addr common.Address, clients []iface.L2Client) (*big.Int, error) {

	if len(clients) < 1 {
		return nil, fmt.Errorf("no client to query sequencer set update time")
	}
	for _, client := range clients {
		// l2 sequencer
		l2Seqencer, err := bindings.NewSequencer(addr, client)
		if err != nil {
			log.Warn("failed to connect to sequencer", "error", err)
			continue
		}
		// get sequencer set update time
		updateTime, err := l2Seqencer.UpdateTime(nil)
		if err != nil {
			log.Warn("failed to get sequencer set update time", "error", err)
			continue
		}
		return updateTime, nil
	}
	return nil, fmt.Errorf("no sequencer set update time found after querying all clients")
}

// query epoch update time from gov contract on l2
func GetEpochUpdateTime(addr common.Address, clients []iface.L2Client) (*big.Int, error) {
	if len(clients) < 1 {
		return nil, fmt.Errorf("no client to query epoch update time")
	}
	for _, client := range clients {
		// l2 gov
		l2Gov, err := bindings.NewGov(addr, client)
		if err != nil {
			log.Warn("failed to connect to gov", "error", err)
			continue
		}
		// get epoch update time
		updateTime, err := l2Gov.RollupEpochUpdateTime(nil)
		if err != nil {
			log.Warn("failed to get epoch update time", "error", err)
			continue
		}
		return updateTime, nil

	}
	return nil, fmt.Errorf("no epoch update time found after querying all clients")

}

func UpdateGasLimit(tx *types.Transaction) (*types.Transaction, error) {
	// add buffer to gas limit (*1.2)
	newGasLimit := tx.Gas() * 12 / 10

	var newTx *types.Transaction
	if tx.Type() == types.LegacyTxType {

		newTx = types.NewTx(&types.LegacyTx{
			Nonce:    tx.Nonce(),
			GasPrice: big.NewInt(tx.GasPrice().Int64()),
			Gas:      newGasLimit,
			To:       tx.To(),
			Value:    tx.Value(),
			Data:     tx.Data(),
		})
	} else if tx.Type() == types.DynamicFeeTxType {
		newTx = types.NewTx(&types.DynamicFeeTx{
			Nonce:     tx.Nonce(),
			GasTipCap: big.NewInt(tx.GasTipCap().Int64()),
			GasFeeCap: big.NewInt(tx.GasFeeCap().Int64()),
			Gas:       newGasLimit,
			To:        tx.To(),
			Value:     tx.Value(),
			Data:      tx.Data(),
		})
	} else if tx.Type() == types.BlobTxType {
		newTx = types.NewTx(&types.BlobTx{
			ChainID:    uint256.MustFromBig(tx.ChainId()),
			Nonce:      tx.Nonce(),
			GasTipCap:  uint256.MustFromBig(big.NewInt(tx.GasTipCap().Int64())),
			GasFeeCap:  uint256.MustFromBig(big.NewInt(tx.GasFeeCap().Int64())),
			Gas:        newGasLimit,
			To:         *tx.To(),
			Value:      uint256.MustFromBig(tx.Value()),
			Data:       tx.Data(),
			BlobFeeCap: uint256.MustFromBig(tx.BlobGasFeeCap()),
			BlobHashes: tx.BlobHashes(),
			Sidecar:    tx.BlobTxSidecar(),
		})

	} else {
		return nil, fmt.Errorf("unknown tx type:%v", tx.Type())
	}
	return newTx, nil
}

// send tx to l1 with business logic check
func (r *Rollup) SendTx(tx *types.Transaction) error {

	// judge tx info is valid
	if tx == nil {
		return errors.New("nil tx")
	}

	err := sendTx(r.L1Client, r.cfg.TxFeeLimit, tx)
	if err != nil {
		return err
	}

	// after send tx
	// add to pending txs
	r.pendingTxs.Add(*tx)

	return nil

}

// send tx to l1 with business logic check
func sendTx(client iface.Client, txFeeLimit uint64, tx *types.Transaction) error {
	// fee limit
	if txFeeLimit > 0 {
		var fee uint64
		// calc tx gas fee
		if tx.Type() == types.BlobTxType {
			// blob fee
			fee = tx.BlobGasFeeCap().Uint64() * tx.BlobGas()
			// tx fee
			fee += tx.GasPrice().Uint64() * tx.Gas()
		} else {
			fee = tx.GasPrice().Uint64() * tx.Gas()
		}

		if fee > txFeeLimit {
			return fmt.Errorf("%v:limit=%v,but got=%v", utils.ErrExceedFeeLimit, txFeeLimit, fee)
		}
	}

	return client.SendTransaction(context.Background(), tx)
}

func (sr *Rollup) ReSubmitTx(resend bool, tx *types.Transaction) (*types.Transaction, error) {
	if tx == nil {
		return nil, errors.New("nil tx")
	}

	// for sign
	opts, err := bind.NewKeyedTransactorWithChainID(sr.privKey, sr.chainId)
	if err != nil {
		return nil, fmt.Errorf("new keyedTransaction with chain id error:%v", err)
	}

	method := "replaced tx"
	if resend {
		method = "resubmitted tx"
	}

	// replaced tx info
	log.Info(method,
		"hash", tx.Hash().String(),
		"gas_fee_cap", tx.GasFeeCap().String(),
		"gas_tip", tx.GasTipCap().String(),
		"blob_fee_cap", tx.BlobGasFeeCap().String(),
		"gas", tx.Gas(),
		"nonce", tx.Nonce(),
	)

	tip, gasFeeCap, blobFeeCap, err := sr.GetGasTipAndCap()
	if err != nil {
		log.Error("get tip and cap", "err", err)
	}
	if !resend {
		// bump tip & feeCap
		bumpedFeeCap := calcThresholdValue(tx.GasFeeCap(), tx.Type() == types.BlobTxType)
		bumpedTip := calcThresholdValue(tx.GasTipCap(), tx.Type() == types.BlobTxType)

		// if bumpedTip > tip
		if bumpedTip.Cmp(tip) > 0 {
			tip = bumpedTip
		}

		if bumpedFeeCap.Cmp(gasFeeCap) > 0 {
			gasFeeCap = bumpedFeeCap
		}

		if tx.Type() == types.BlobTxType {
			bumpedBlobFeeCap := calcThresholdValue(tx.BlobGasFeeCap(), tx.Type() == types.BlobTxType)
			if bumpedBlobFeeCap.Cmp(blobFeeCap) > 0 {
				blobFeeCap = bumpedBlobFeeCap
			}
		}
	}

	var newTx *types.Transaction
	switch tx.Type() {
	case types.DynamicFeeTxType:
		newTx = types.NewTx(&types.DynamicFeeTx{
			To:        tx.To(),
			Nonce:     tx.Nonce(),
			GasFeeCap: gasFeeCap,
			GasTipCap: tip,
			Gas:       tx.Gas(),
			Value:     tx.Value(),
			Data:      tx.Data(),
		})
	case types.BlobTxType:

		newTx = types.NewTx(&types.BlobTx{
			ChainID:    uint256.MustFromBig(tx.ChainId()),
			Nonce:      tx.Nonce(),
			GasTipCap:  uint256.MustFromBig(tip),
			GasFeeCap:  uint256.MustFromBig(gasFeeCap),
			Gas:        tx.Gas(),
			To:         *tx.To(),
			Value:      uint256.MustFromBig(tx.Value()),
			Data:       tx.Data(),
			BlobFeeCap: uint256.MustFromBig(blobFeeCap),
			BlobHashes: tx.BlobHashes(),
			Sidecar:    tx.BlobTxSidecar(),
		})

	default:
		return nil, fmt.Errorf("replace unknown tx type:%v", tx.Type())

	}

	log.Info("new tx info",
		"tx_type", newTx.Type(),
		"gas_tip", tip.String(), //todo: convert to gwei
		"gas_fee_cap", gasFeeCap.String(), //todo: convert to gwei
		"blob_fee_cap", blobFeeCap.String(), //todo: convert to gwei
	)
	// sign tx
	opts.Nonce = big.NewInt(int64(newTx.Nonce()))
	newTx, err = opts.Signer(opts.From, newTx)
	if err != nil {
		return nil, fmt.Errorf("sign tx error:%w", err)
	}
	// send tx
	err = sr.SendTx(newTx)
	if err != nil {
		return nil, fmt.Errorf("send tx error:%w", err)
	}

	return newTx, nil
}

func (r *Rollup) IsStaker() (bool, error) {

	isStaker, err := r.Staking.IsStaker(nil, common.HexToAddress(r.walletAddr()))
	if err != nil {
		return false, fmt.Errorf("call IsStaker err :%v", err)
	}
	return isStaker, nil
}

func (sr *Rollup) EstimateGas(from, to common.Address, data []byte, feecap *big.Int, tip *big.Int) (uint64, error) {

	gas, err := sr.L1Client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:      from,
		To:        &to,
		GasFeeCap: feecap,
		GasTipCap: tip,
		Data:      data,
	})
	return gas, err

}

// for rollup
func (r *Rollup) RoughEstimateGas(msgcnt uint64) uint64 {
	return r.cfg.RollupTxGasBase + msgcnt*r.cfg.RollupTxGasPerL1Msg
}
