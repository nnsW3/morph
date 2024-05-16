// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"encoding/json"

	"morph-l2/bindings/solc"
)

const L2WETHGatewayStorageLayoutJSON = "{\"storage\":[{\"astId\":1000,\"contract\":\"contracts/l2/gateways/L2WETHGateway.sol:L2WETHGateway\",\"label\":\"_initialized\",\"offset\":0,\"slot\":\"0\",\"type\":\"t_uint8\"},{\"astId\":1001,\"contract\":\"contracts/l2/gateways/L2WETHGateway.sol:L2WETHGateway\",\"label\":\"_initializing\",\"offset\":1,\"slot\":\"0\",\"type\":\"t_bool\"},{\"astId\":1002,\"contract\":\"contracts/l2/gateways/L2WETHGateway.sol:L2WETHGateway\",\"label\":\"_status\",\"offset\":0,\"slot\":\"1\",\"type\":\"t_uint256\"},{\"astId\":1003,\"contract\":\"contracts/l2/gateways/L2WETHGateway.sol:L2WETHGateway\",\"label\":\"__gap\",\"offset\":0,\"slot\":\"2\",\"type\":\"t_array(t_uint256)1013_storage\"},{\"astId\":1004,\"contract\":\"contracts/l2/gateways/L2WETHGateway.sol:L2WETHGateway\",\"label\":\"__gap\",\"offset\":0,\"slot\":\"51\",\"type\":\"t_array(t_uint256)1014_storage\"},{\"astId\":1005,\"contract\":\"contracts/l2/gateways/L2WETHGateway.sol:L2WETHGateway\",\"label\":\"_owner\",\"offset\":0,\"slot\":\"101\",\"type\":\"t_address\"},{\"astId\":1006,\"contract\":\"contracts/l2/gateways/L2WETHGateway.sol:L2WETHGateway\",\"label\":\"__gap\",\"offset\":0,\"slot\":\"102\",\"type\":\"t_array(t_uint256)1013_storage\"},{\"astId\":1007,\"contract\":\"contracts/l2/gateways/L2WETHGateway.sol:L2WETHGateway\",\"label\":\"counterpart\",\"offset\":0,\"slot\":\"151\",\"type\":\"t_address\"},{\"astId\":1008,\"contract\":\"contracts/l2/gateways/L2WETHGateway.sol:L2WETHGateway\",\"label\":\"router\",\"offset\":0,\"slot\":\"152\",\"type\":\"t_address\"},{\"astId\":1009,\"contract\":\"contracts/l2/gateways/L2WETHGateway.sol:L2WETHGateway\",\"label\":\"messenger\",\"offset\":0,\"slot\":\"153\",\"type\":\"t_address\"},{\"astId\":1010,\"contract\":\"contracts/l2/gateways/L2WETHGateway.sol:L2WETHGateway\",\"label\":\"__gap\",\"offset\":0,\"slot\":\"154\",\"type\":\"t_array(t_uint256)1012_storage\"},{\"astId\":1011,\"contract\":\"contracts/l2/gateways/L2WETHGateway.sol:L2WETHGateway\",\"label\":\"__gap\",\"offset\":0,\"slot\":\"200\",\"type\":\"t_array(t_uint256)1014_storage\"}],\"types\":{\"t_address\":{\"encoding\":\"inplace\",\"label\":\"address\",\"numberOfBytes\":\"20\"},\"t_array(t_uint256)1012_storage\":{\"encoding\":\"inplace\",\"label\":\"uint256[46]\",\"numberOfBytes\":\"1472\"},\"t_array(t_uint256)1013_storage\":{\"encoding\":\"inplace\",\"label\":\"uint256[49]\",\"numberOfBytes\":\"1568\"},\"t_array(t_uint256)1014_storage\":{\"encoding\":\"inplace\",\"label\":\"uint256[50]\",\"numberOfBytes\":\"1600\"},\"t_bool\":{\"encoding\":\"inplace\",\"label\":\"bool\",\"numberOfBytes\":\"1\"},\"t_uint256\":{\"encoding\":\"inplace\",\"label\":\"uint256\",\"numberOfBytes\":\"32\"},\"t_uint8\":{\"encoding\":\"inplace\",\"label\":\"uint8\",\"numberOfBytes\":\"1\"}}}"

var L2WETHGatewayStorageLayout = new(solc.StorageLayout)

var L2WETHGatewayDeployedBin = "0x6080604052600436106100e7575f3560e01c8063797594b011610087578063c0c53b8b11610057578063c0c53b8b14610346578063c676ad2914610365578063f2fde38b146103a4578063f887ea40146103c3575f80fd5b8063797594b0146102ca5780638431f5c1146102f65780638da5cb5b14610309578063a93a4af914610333575f80fd5b806354bbd59c116100c257806354bbd59c14610251578063575361b6146102905780636c07ea43146102a3578063715018a6146102b6575f80fd5b806319c4d4c6146101965780631efd482a146101f25780633cb747bf14610225575f80fd5b3661019257337f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1614610190576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6f6e6c792057455448000000000000000000000000000000000000000000000060448201526064015b60405180910390fd5b005b5f80fd5b3480156101a1575f80fd5b506101c97f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200160405180910390f35b3480156101fd575f80fd5b506101c97f000000000000000000000000000000000000000000000000000000000000000081565b348015610230575f80fd5b506099546101c99073ffffffffffffffffffffffffffffffffffffffff1681565b34801561025c575f80fd5b506101c961026b366004611b2f565b507f000000000000000000000000000000000000000000000000000000000000000090565b61019061029e366004611b96565b6103ef565b6101906102b1366004611c0c565b61043a565b3480156102c1575f80fd5b50610190610478565b3480156102d5575f80fd5b506097546101c99073ffffffffffffffffffffffffffffffffffffffff1681565b610190610304366004611c3e565b61048b565b348015610314575f80fd5b5060655473ffffffffffffffffffffffffffffffffffffffff166101c9565b610190610341366004611cd0565b61092d565b348015610351575f80fd5b50610190610360366004611d13565b61093f565b348015610370575f80fd5b506101c961037f366004611b2f565b507f000000000000000000000000000000000000000000000000000000000000000090565b3480156103af575f80fd5b506101906103be366004611b2f565b610b4d565b3480156103ce575f80fd5b506098546101c99073ffffffffffffffffffffffffffffffffffffffff1681565b61043286868686868080601f0160208091040260200160405190810160405280939291908181526020018383808284375f92019190915250889250610c04915050565b505050505050565b6104738333845f5b6040519080825280601f01601f19166020018201604052801561046c576020820181803683370190505b5085610c04565b505050565b610480611090565b6104895f611111565b565b60995473ffffffffffffffffffffffffffffffffffffffff1633811461050d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f6f6e6c79206d657373656e6765722063616e2063616c6c0000000000000000006044820152606401610187565b8073ffffffffffffffffffffffffffffffffffffffff16636e296e456040518163ffffffff1660e01b8152600401602060405180830381865afa158015610556573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061057a9190611d88565b60975473ffffffffffffffffffffffffffffffffffffffff9081169116146105fe576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f6f6e6c792063616c6c20627920636f756e7465727061727400000000000000006044820152606401610187565b610606611187565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168873ffffffffffffffffffffffffffffffffffffffff16146106bb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f6c3120746f6b656e206e6f7420574554480000000000000000000000000000006044820152606401610187565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168773ffffffffffffffffffffffffffffffffffffffff1614610770576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f6c3220746f6b656e206e6f7420574554480000000000000000000000000000006044820152606401610187565b3484146107d9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f6d73672e76616c7565206d69736d6174636800000000000000000000000000006044820152606401610187565b8673ffffffffffffffffffffffffffffffffffffffff1663d0e30db0856040518263ffffffff1660e01b81526004015f604051808303818588803b15801561081f575f80fd5b505af1158015610831573d5f803e3d5ffd5b506108599350505073ffffffffffffffffffffffffffffffffffffffff8916905086866111fa565b6108988584848080601f0160208091040260200160405190810160405280939291908181526020018383808284375f920191909152506112ce92505050565b8573ffffffffffffffffffffffffffffffffffffffff168773ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff167f165ba69f6ab40c50cade6f65431801e5f9c7d7830b7545391920db039133ba34888888886040516109129493929190611da3565b60405180910390a461092360018055565b5050505050505050565b6109398484845f610442565b50505050565b5f54610100900460ff161580801561095d57505f54600160ff909116105b806109765750303b15801561097657505f5460ff166001145b610a02576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610187565b5f80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790558015610a5e575f80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b73ffffffffffffffffffffffffffffffffffffffff8316610adb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f7a65726f20726f757465722061646472657373000000000000000000000000006044820152606401610187565b610ae684848461137e565b8015610939575f80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150505050565b610b55611090565b73ffffffffffffffffffffffffffffffffffffffff8116610bf8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610187565b610c0181611111565b50565b610c0c611187565b5f8311610c75576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f7769746864726177207a65726f20616d6f756e740000000000000000000000006044820152606401610187565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1614610d2a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f6f6e6c79205745544820697320616c6c6f7765640000000000000000000000006044820152606401610187565b609854339073ffffffffffffffffffffffffffffffffffffffff16819003610d655782806020019051810190610d609190611e35565b935090505b610d8773ffffffffffffffffffffffffffffffffffffffff8716823087611529565b6040517f2e1a7d4d0000000000000000000000000000000000000000000000000000000081526004810185905273ffffffffffffffffffffffffffffffffffffffff871690632e1a7d4d906024015f604051808303815f87803b158015610dec575f80fd5b505af1158015610dfe573d5f803e3d5ffd5b50506040517f000000000000000000000000000000000000000000000000000000000000000092505f9150610e419083908a9086908b908b908b90602401611f5a565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152602080830180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167f84bd13b00000000000000000000000000000000000000000000000000000000017905260995482517fecc7042800000000000000000000000000000000000000000000000000000000815292519394505f9373ffffffffffffffffffffffffffffffffffffffff9091169263ecc704289260048083019391928290030181865afa158015610f25573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610f499190611fb4565b60995490915073ffffffffffffffffffffffffffffffffffffffff1663b2267a7b610f74348a611fcb565b6097546040517fffffffff0000000000000000000000000000000000000000000000000000000060e085901b168152610fcc9173ffffffffffffffffffffffffffffffffffffffff16908c9088908c90600401612009565b5f604051808303818588803b158015610fe3575f80fd5b505af1158015610ff5573d5f803e3d5ffd5b50505050508373ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fa9967b740f3fc69dfbf4744b4b1c7dfdb0b1b63f1fa4cf573bcdcb9f3ac687c48b8b8b876040516110749493929190612009565b60405180910390a45050505061108960018055565b5050505050565b60655473ffffffffffffffffffffffffffffffffffffffff163314610489576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610187565b6065805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a35050565b6002600154036111f3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c006044820152606401610187565b6002600155565b60405173ffffffffffffffffffffffffffffffffffffffff83166024820152604481018290526104739084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152611587565b5f81511180156112f457505f8273ffffffffffffffffffffffffffffffffffffffff163b115b15611374576040517f444b281f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83169063444b281f9061134b90849060040161204e565b5f604051808303815f87803b158015611362575f80fd5b505af1158015610432573d5f803e3d5ffd5b5050565b60018055565b73ffffffffffffffffffffffffffffffffffffffff83166113fb576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f7a65726f20636f756e74657270617274206164647265737300000000000000006044820152606401610187565b73ffffffffffffffffffffffffffffffffffffffff8116611478576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f7a65726f206d657373656e6765722061646472657373000000000000000000006044820152606401610187565b611480611694565b611488611732565b6097805473ffffffffffffffffffffffffffffffffffffffff8086167fffffffffffffffffffffffff000000000000000000000000000000000000000092831617909255609980548484169216919091179055821615610473576098805473ffffffffffffffffffffffffffffffffffffffff84167fffffffffffffffffffffffff0000000000000000000000000000000000000000909116179055505050565b60405173ffffffffffffffffffffffffffffffffffffffff808516602483015283166044820152606481018290526109399085907f23b872dd000000000000000000000000000000000000000000000000000000009060840161124c565b5f6115e8826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166117d09092919063ffffffff16565b905080515f14806116085750808060200190518101906116089190612060565b610473576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610187565b5f54610100900460ff1661172a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610187565b6104896117e6565b5f54610100900460ff166117c8576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610187565b61048961187c565b60606117de84845f8561191b565b949350505050565b5f54610100900460ff16611378576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610187565b5f54610100900460ff16611912576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610187565b61048933611111565b6060824710156119ad576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c00000000000000000000000000000000000000000000000000006064820152608401610187565b5f808673ffffffffffffffffffffffffffffffffffffffff1685876040516119d5919061207f565b5f6040518083038185875af1925050503d805f8114611a0f576040519150601f19603f3d011682016040523d82523d5f602084013e611a14565b606091505b5091509150611a2587838387611a30565b979650505050505050565b60608315611ac55782515f03611abe5773ffffffffffffffffffffffffffffffffffffffff85163b611abe576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610187565b50816117de565b6117de8383815115611ada5781518083602001fd5b806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610187919061204e565b73ffffffffffffffffffffffffffffffffffffffff81168114610c01575f80fd5b5f60208284031215611b3f575f80fd5b8135611b4a81611b0e565b9392505050565b5f8083601f840112611b61575f80fd5b50813567ffffffffffffffff811115611b78575f80fd5b602083019150836020828501011115611b8f575f80fd5b9250929050565b5f805f805f8060a08789031215611bab575f80fd5b8635611bb681611b0e565b95506020870135611bc681611b0e565b945060408701359350606087013567ffffffffffffffff811115611be8575f80fd5b611bf489828a01611b51565b979a9699509497949695608090950135949350505050565b5f805f60608486031215611c1e575f80fd5b8335611c2981611b0e565b95602085013595506040909401359392505050565b5f805f805f805f60c0888a031215611c54575f80fd5b8735611c5f81611b0e565b96506020880135611c6f81611b0e565b95506040880135611c7f81611b0e565b94506060880135611c8f81611b0e565b93506080880135925060a088013567ffffffffffffffff811115611cb1575f80fd5b611cbd8a828b01611b51565b989b979a50959850939692959293505050565b5f805f8060808587031215611ce3575f80fd5b8435611cee81611b0e565b93506020850135611cfe81611b0e565b93969395505050506040820135916060013590565b5f805f60608486031215611d25575f80fd5b8335611d3081611b0e565b92506020840135611d4081611b0e565b91506040840135611d5081611b0e565b809150509250925092565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f60208284031215611d98575f80fd5b8151611b4a81611b0e565b73ffffffffffffffffffffffffffffffffffffffff8516815283602082015260606040820152816060820152818360808301375f818301608090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01601019392505050565b5f5b83811015611e2d578181015183820152602001611e15565b50505f910152565b5f8060408385031215611e46575f80fd5b8251611e5181611b0e565b602084015190925067ffffffffffffffff80821115611e6e575f80fd5b818501915085601f830112611e81575f80fd5b815181811115611e9357611e93611d5b565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f01168101908382118183101715611ed957611ed9611d5b565b81604052828152886020848701011115611ef1575f80fd5b611f02836020830160208801611e13565b80955050505050509250929050565b5f8151808452611f28816020860160208601611e13565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b5f73ffffffffffffffffffffffffffffffffffffffff80891683528088166020840152808716604084015280861660608401525083608083015260c060a0830152611fa860c0830184611f11565b98975050505050505050565b5f60208284031215611fc4575f80fd5b5051919050565b80820180821115612003577f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b92915050565b73ffffffffffffffffffffffffffffffffffffffff85168152836020820152608060408201525f61203d6080830185611f11565b905082606083015295945050505050565b602081525f611b4a6020830184611f11565b5f60208284031215612070575f80fd5b81518015158114611b4a575f80fd5b5f8251612090818460208701611e13565b919091019291505056fea164736f6c6343000818000a"

func init() {
	if err := json.Unmarshal([]byte(L2WETHGatewayStorageLayoutJSON), L2WETHGatewayStorageLayout); err != nil {
		panic(err)
	}

	layouts["L2WETHGateway"] = L2WETHGatewayStorageLayout
	deployedBytecodes["L2WETHGateway"] = L2WETHGatewayDeployedBin
}
