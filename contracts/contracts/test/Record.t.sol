// SPDX-License-Identifier: MIT
pragma solidity =0.8.24;

import {L2StakingBaseTest} from "./base/L2StakingBase.t.sol";
import {IRecord} from "../l2/staking/IRecord.sol";

contract RecordTest is L2StakingBaseTest {
    function setUp() public virtual override {
        super.setUp();
    }

    /**
     * @notice initialize: re-initialize
     */
    function test_initialize_onlyOnce_reverts() public {
        hevm.expectRevert("Initializable: contract is already initialized");
        hevm.prank(multisig);
        record.initialize(address(0), 1);
    }

    /**
     * @notice setOracleAddress: check params
     */
    function test_setOracleAddress_invalidAddress_reverts() public {
        hevm.expectRevert("invalid oracle address");
        hevm.prank(multisig);
        record.setOracleAddress(address(0));
    }

    /**
     * @notice setOracleAddress: check owner
     */
    function test_owner_onlyOwner_reverts() public {
        hevm.expectRevert("Ownable: caller is not the owner");
        hevm.prank(alice);
        record.setOracleAddress(address(0));
    }

    /**
     * @notice setLatestRewardEpochBlock: check params
     */
    function test_setLatestRewardEpochBlock_paramsCheck_reverts() public {
        hevm.expectRevert("only oracle allowed");
        hevm.prank(multisig);
        record.setLatestRewardEpochBlock(0);

        hevm.expectRevert("invalid latest block");
        hevm.prank(oracleAddress);
        record.setLatestRewardEpochBlock(0);

        hevm.prank(oracleAddress);
        record.setLatestRewardEpochBlock(100);

        hevm.expectRevert("already set");
        hevm.prank(oracleAddress);
        record.setLatestRewardEpochBlock(100);
    }

    /**
     * @notice recordFinalizedBatchSubmissions
     * 1. check owner
     * 2. check params
     */
    function test_recordFinalizedBatchSubmissions_paramsCheck_reverts() public {
        IRecord.BatchSubmission[] memory submissions = new IRecord.BatchSubmission[](0);

        hevm.expectRevert("only oracle allowed");
        hevm.prank(multisig);
        record.recordFinalizedBatchSubmissions(submissions);

        hevm.expectRevert("empty batch submissions");
        hevm.prank(oracleAddress);
        record.recordFinalizedBatchSubmissions(submissions);

        submissions = new IRecord.BatchSubmission[](1);
        hevm.expectRevert("invalid index");
        hevm.prank(oracleAddress);
        record.recordFinalizedBatchSubmissions(submissions);

        // recordFinalizedBatchSubmissions
        IRecord.BatchSubmission memory submission = IRecord.BatchSubmission(
            nextBatchSubmissionIndex,
            address(0),
            0,
            1,
            0,
            0
        );
        submissions = new IRecord.BatchSubmission[](1);
        submissions[0] = submission;
        hevm.expectEmit(true, true, false, false);
        emit IRecord.BatchSubmissionsUploaded(1, 1);
        hevm.prank(oracleAddress);
        record.recordFinalizedBatchSubmissions(submissions);
    }

    /**
     * @notice recordRollupEpochs
     * 1. check owner
     * 2. check params
     */
    function test_recordRollupEpochs_paramsCheck_reverts() public {
        IRecord.RollupEpochInfo[] memory epochInfos = new IRecord.RollupEpochInfo[](0);

        hevm.expectRevert("only oracle allowed");
        hevm.prank(multisig);
        record.recordRollupEpochs(epochInfos);

        hevm.expectRevert("empty rollup epochs");
        hevm.prank(oracleAddress);
        record.recordRollupEpochs(epochInfos);

        epochInfos = new IRecord.RollupEpochInfo[](1);
        IRecord.RollupEpochInfo memory epochInfo = IRecord.RollupEpochInfo(
            1, // invalid index
            address(0),
            0,
            0,
            0
        );
        epochInfos[0] = epochInfo;
        hevm.expectRevert("invalid index");
        hevm.prank(oracleAddress);
        record.recordRollupEpochs(epochInfos);

        // recordRollupEpochs
        epochInfo = IRecord.RollupEpochInfo(0, address(0), 0, 0, 0);
        epochInfos[0] = epochInfo;
        hevm.expectEmit(true, true, false, false);
        emit IRecord.RollupEpochsUploaded(0, 1);
        hevm.prank(oracleAddress);
        record.recordRollupEpochs(epochInfos);
    }

    /**
     * @notice recordRewardEpochs: check owner
     */
    function test_recordRewardEpochs_onlyOwner_reverts() public {
        uint256 sequencerSize = sequencer.getSequencerSet2Size();
        address[] memory sequencers = sequencer.getSequencerSet2();
        uint256[] memory sequencerBlocks = new uint256[](sequencerSize);
        uint256[] memory sequencerRatios = new uint256[](sequencerSize);
        uint256[] memory sequencerCommissions = new uint256[](sequencerSize);

        for (uint256 i = 0; i < sequencerSize; i++) {
            sequencerBlocks[i] = 0;
            sequencerRatios[i] = SEQUENCER_RATIO_PRECISION / sequencerSize;
            sequencerCommissions[i] = 1;
        }

        IRecord.RewardEpochInfo[] memory rewardEpochInfos = new IRecord.RewardEpochInfo[](1);

        rewardEpochInfos[0] = IRecord.RewardEpochInfo(
            0,
            1,
            sequencers,
            sequencerBlocks,
            sequencerRatios,
            sequencerCommissions
        );

        hevm.expectRevert("only oracle allowed");
        hevm.prank(multisig);
        record.recordRewardEpochs(rewardEpochInfos);
    }

    /**
     * @notice recordRewardEpochs: check params
     */
    function test_recordRewardEpochs_paramsCheck_reverts() public {
        uint256 sequencerSize = sequencer.getSequencerSet2Size();
        address[] memory sequencers = sequencer.getSequencerSet2();
        uint256[] memory sequencerBlocks = new uint256[](sequencerSize);
        uint256[] memory sequencerRatios = new uint256[](sequencerSize);
        uint256[] memory sequencerCommissions = new uint256[](sequencerSize);

        for (uint256 i = 0; i < sequencerSize; i++) {
            sequencerBlocks[i] = 1;
            sequencerRatios[i] = SEQUENCER_RATIO_PRECISION / sequencerSize;
            sequencerCommissions[i] = 1;
        }

        IRecord.RewardEpochInfo[] memory rewardEpochInfos = new IRecord.RewardEpochInfo[](0);

        hevm.expectRevert("empty reward epochs");
        hevm.prank(oracleAddress);
        record.recordRewardEpochs(rewardEpochInfos);

        // greater than minted epoch
        rewardEpochInfos = new IRecord.RewardEpochInfo[](2);

        hevm.expectRevert("start block should be set");
        hevm.prank(oracleAddress);
        record.recordRewardEpochs(rewardEpochInfos);

        hevm.prank(oracleAddress);
        record.setLatestRewardEpochBlock(1);

        hevm.expectRevert("reward is not started yet");
        hevm.prank(oracleAddress);
        record.recordRewardEpochs(rewardEpochInfos);

        // update epoch
        hevm.warp(rewardStartTime * 2);
        rewardEpochInfos = new IRecord.RewardEpochInfo[](1);
        rewardEpochInfos[0] = IRecord.RewardEpochInfo(
            0,
            1, // total block not equal
            sequencers,
            sequencerBlocks,
            sequencerRatios,
            sequencerCommissions
        );
        hevm.expectRevert("invalid sequencers blocks");
        hevm.prank(oracleAddress);
        record.recordRewardEpochs(rewardEpochInfos);

        // invalide commission rate
        sequencerCommissions = new uint256[](sequencerSize);
        for (uint256 i = 0; i < sequencerSize; i++) {
            sequencerCommissions[i] = 21;
        }
        rewardEpochInfos[0] = IRecord.RewardEpochInfo(
            0,
            1, // total block not equal
            sequencers,
            sequencerBlocks,
            sequencerRatios,
            sequencerCommissions
        );

        hevm.expectRevert("invalid sequencers commission");
        hevm.prank(oracleAddress);
        record.recordRewardEpochs(rewardEpochInfos);

        // invalide sequencers ratios
        sequencerRatios = new uint256[](sequencerSize);
        sequencerCommissions = new uint256[](sequencerSize);
        for (uint256 i = 0; i < sequencerSize; i++) {
            sequencerRatios[i] = SEQUENCER_RATIO_PRECISION / sequencerSize + 1;
            sequencerCommissions[i] = 2;
        }
        rewardEpochInfos[0] = IRecord.RewardEpochInfo(
            0,
            3, // total block not equal
            sequencers,
            sequencerBlocks,
            sequencerRatios,
            sequencerCommissions
        );

        hevm.expectRevert("invalid sequencers ratios");
        hevm.prank(oracleAddress);
        record.recordRewardEpochs(rewardEpochInfos);
    }

    /**
     * @notice getBatchSubmissions: check params
     */
    function test_getBatchSubmissions_paramsCheck_reverts() public {
        hevm.expectRevert("invalid index");
        hevm.prank(oracleAddress);
        record.getBatchSubmissions(2, 1);
    }

    /**
     * @notice getRollupEpochs: check params
     */
    function test_getRollupEpochs_paramsCheck_reverts() public {
        hevm.expectRevert("invalid index");
        hevm.prank(oracleAddress);
        record.getRollupEpochs(2, 1);
    }

    /**
     * @notice getRewardEpochs: check params
     */
    function test_getRewardEpochs_paramsCheck_reverts() public {
        hevm.expectRevert("invalid index");
        hevm.prank(oracleAddress);
        record.getRewardEpochs(2, 1);
    }
}
