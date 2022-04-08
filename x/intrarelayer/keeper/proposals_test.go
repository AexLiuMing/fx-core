package keeper_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/functionx/fx-core/tests"
	fxcoretypes "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/intrarelayer/types"
	"github.com/functionx/fx-core/x/intrarelayer/types/contracts"
	abci "github.com/tendermint/tendermint/abci/types"
	"math/big"
	"strings"
)

const (
	fip20Name       = "coin"
	fip20Symbol     = "token"
	cosmosTokenName = "coin"
	displayCoinName = "COIN"
	defaultExponent = uint32(18)
	zeroExponent    = uint32(0)
)

func (suite *KeeperTestSuite) setupRegisterFIP20Pair() common.Address {
	suite.SetupTest()
	contractAddr := suite.DeployContract(suite.address, fip20Name, fip20Symbol, 18)
	suite.Commit()
	_, err := suite.app.IntrarelayerKeeper.RegisterFIP20(suite.ctx, contractAddr)
	suite.Require().NoError(err)
	return contractAddr
}

func (suite *KeeperTestSuite) setupRegisterCoin() (banktypes.Metadata, *types.TokenPair) {
	suite.SetupTest()
	validMetadata := banktypes.Metadata{
		Description: "desc",
		Base:        cosmosTokenName,
		// NOTE: Denom units MUST be increasing
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    cosmosTokenName,
				Exponent: 0,
			},
			{
				Denom:    displayCoinName,
				Exponent: uint32(18),
			},
		},
		Display: displayCoinName,
	}
	// pair := types.NewTokenPair(contractAddr, cosmosTokenName, true, types.OWNER_MODULE)
	pair, err := suite.app.IntrarelayerKeeper.RegisterCoin(suite.ctx, validMetadata)
	suite.Require().NoError(err)
	suite.Commit()
	return validMetadata, pair
}

func (suite KeeperTestSuite) TestRegisterCoin() {
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"intrarelaying is disabled globally",
			func() {
				params := types.DefaultParams()
				params.EnableIntrarelayer = false
				suite.app.IntrarelayerKeeper.SetParams(suite.ctx, params)
			},
			false,
		},
		{
			"denom already registered",
			func() {
				regPair := types.NewTokenPair(tests.GenerateAddress(), cosmosTokenName, true, types.OWNER_MODULE)
				suite.app.IntrarelayerKeeper.SetDenomMap(suite.ctx, regPair.Denom, regPair.GetID())
				suite.Commit()
			},
			false,
		},
		{
			"metadata different that stored",
			func() {
				validMetadata := banktypes.Metadata{
					Description: "desc",
					Base:        cosmosTokenName,
					// NOTE: Denom units MUST be increasing
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    cosmosTokenName,
							Exponent: 0,
						},
						{
							Denom:    displayCoinName,
							Exponent: uint32(1),
						},
						{
							Denom:    "extraDenom",
							Exponent: uint32(2),
						},
					},
					//Name:    "otherName",
					//Symbol:  "token",
					Display: displayCoinName,
				}
				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, validMetadata)
			},
			false,
		},
		{
			"ok",
			func() {},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()
			validMetadata := banktypes.Metadata{
				Description: "desc",
				Base:        cosmosTokenName,
				// NOTE: Denom units MUST be increasing
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    cosmosTokenName,
						Exponent: 0,
					},
					{
						Denom:    displayCoinName,
						Exponent: uint32(1),
					},
				},
				Display: displayCoinName,
			}
			pair, err := suite.app.IntrarelayerKeeper.RegisterCoin(suite.ctx, validMetadata)
			suite.Commit()
			expPair := &types.TokenPair{
				Fip20Address:  "0x00819E780C6e96c50Ed70eFFf5B73569c15d0bd7",
				Denom:         "coin",
				Enabled:       true,
				ContractOwner: 1,
			}
			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				suite.Require().Equal(pair, expPair)
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

func (suite KeeperTestSuite) TestUpgradeContract() {
	suite.SetupTest()
	acc := suite.app.EvmKeeper.GetAccount(suite.ctx, common.HexToAddress(contracts.FIP20UpgradeCodeAddress))
	suite.Require().NotEmpty(acc)
	suite.Require().True(acc.IsContract())
	code := suite.app.EvmKeeper.GetCode(suite.ctx, common.BytesToHash(acc.CodeHash))
	suite.Require().True(len(code) > 0)
}

func (suite KeeperTestSuite) TestDeployTokenUpgrade() {
	if fxcoretypes.Network() != fxcoretypes.NetworkDevnet() {
		suite.T().SkipNow()
	}
	suite.SetupTest()

	addr, err := suite.app.IntrarelayerKeeper.DeployTokenUpgrade(suite.ctx, types.ModuleAddress, "Function X Token", "FX", 18, false)
	suite.Require().NoError(err)
	suite.T().Log("addr", addr.String())

	fip20ABI := contracts.MustGetABI(suite.ctx.BlockHeight(), contracts.FIP20UpgradeType)
	_, err = suite.app.IntrarelayerKeeper.CallEVM(suite.ctx, fip20ABI, types.ModuleAddress, addr, "mint", suite.address, big.NewInt(10000))
	suite.Require().NoError(err)

	balance, err := suite.app.IntrarelayerKeeper.QueryFIP20BalanceOf(suite.ctx, addr, suite.address)
	suite.Require().NoError(err)
	suite.T().Log("balance", balance.String())

	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 9)
	suite.app.IntrarelayerKeeper.BeginBlock(suite.ctx, abci.RequestBeginBlock{})

	_, err = suite.app.IntrarelayerKeeper.CallEVM(suite.ctx, fip20ABI, types.ModuleAddress, addr, "mint", suite.address, big.NewInt(10000))
	suite.Require().NoError(err)

	balance, err = suite.app.IntrarelayerKeeper.QueryFIP20BalanceOf(suite.ctx, addr, suite.address)
	suite.Require().NoError(err)
	suite.T().Log("balance", balance.String())

	metadata, err := suite.app.IntrarelayerKeeper.QueryFIP20(suite.ctx, addr)
	suite.Require().NoError(err)
	suite.T().Log(metadata.Name, metadata.Symbol, metadata.Decimals)
}

func (suite KeeperTestSuite) TestDeployContractFIP20Upgrade() {
	suite.SetupTest()

	fip20ABI := contracts.MustGetABI(suite.ctx.BlockHeight(), contracts.FIP20UpgradeType)
	fip20Bin := contracts.MustGetBin(suite.ctx.BlockHeight(), contracts.FIP20UpgradeType)

	addr, err := suite.app.IntrarelayerKeeper.DeployContract(suite.ctx, types.ModuleAddress, fip20ABI, fip20Bin)
	suite.Require().NoError(err)

	account := suite.app.EvmKeeper.GetAccount(suite.ctx, addr)
	suite.Require().NotEmpty(account)

	code := suite.app.EvmKeeper.GetCode(suite.ctx, common.BytesToHash(account.CodeHash))
	newCode := bytes.ReplaceAll(code, addr.Bytes(), common.HexToAddress(contracts.EmptyAddress).Bytes())

	cCode := contracts.MustGetCode(suite.ctx.BlockHeight(), contracts.FIP20UpgradeType)
	suite.Require().Equal(cCode, newCode)
}

func (suite KeeperTestSuite) TestDeployContractFIP20() {
	suite.SetupTest()

	abiJson, _ := abi.JSON(strings.NewReader("[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"decimals_\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"to\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"target\",\"type\":\"string\"}],\"name\":\"TransferCross\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isContract\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"module\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"to\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"target\",\"type\":\"string\"}],\"name\":\"transferCross\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"))
	cBin, _ := hex.DecodeString("60806040523480156200001157600080fd5b50604051620010fc380380620010fc8339810160408190526200003491620001e8565b8251620000499060009060208601906200008f565b5081516200005f9060019060208501906200008f565b506002805460ff90921660ff199092169190911790555050600680546001600160a01b03191633179055620002bc565b8280546200009d9062000269565b90600052602060002090601f016020900481019282620000c157600085556200010c565b82601f10620000dc57805160ff19168380011785556200010c565b828001600101855582156200010c579182015b828111156200010c578251825591602001919060010190620000ef565b506200011a9291506200011e565b5090565b5b808211156200011a57600081556001016200011f565b600082601f83011262000146578081fd5b81516001600160401b0380821115620001635762000163620002a6565b604051601f8301601f19908116603f011681019082821181831017156200018e576200018e620002a6565b81604052838152602092508683858801011115620001aa578485fd5b8491505b83821015620001cd5785820183015181830184015290820190620001ae565b83821115620001de57848385830101525b9695505050505050565b600080600060608486031215620001fd578283fd5b83516001600160401b038082111562000214578485fd5b620002228783880162000135565b9450602086015191508082111562000238578384fd5b50620002478682870162000135565b925050604084015160ff811681146200025e578182fd5b809150509250925092565b6002810460018216806200027e57607f821691505b60208210811415620002a057634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052604160045260246000fd5b610e3080620002cc6000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c806370a0823111610097578063a2ac8a3711610066578063a2ac8a3714610212578063a9059cbb14610225578063b86d529814610238578063dd62ed3e14610249576100f5565b806370a08231146101ac5780638da5cb5b146101cc57806395d89b41146101f75780639dc29fac146101ff576100f5565b806318160ddd116100d357806318160ddd1461014e57806323b872dd14610165578063313ce5671461017857806340c10f1914610197576100f5565b806306fdde03146100fa578063095ea7b314610118578063162790551461013b575b600080fd5b610102610274565b60405161010f9190610d15565b60405180910390f35b61012b610126366004610c2d565b610302565b604051901515815260200161010f565b61012b610149366004610b9f565b610358565b61015760035481565b60405190815260200161010f565b61012b610173366004610bf2565b610368565b6002546101859060ff1681565b60405160ff909116815260200161010f565b6101aa6101a5366004610c2d565b610415565b005b6101576101ba366004610b9f565b60046020526000908152604090205481565b6006546101df906001600160a01b031681565b6040516001600160a01b03909116815260200161010f565b610102610477565b6101aa61020d366004610c2d565b610484565b61012b610220366004610c56565b6104e2565b61012b610233366004610c2d565b610547565b6006546001600160a01b03166101df565b610157610257366004610bc0565b600560209081526000928352604080842090915290825290205481565b6000805461028190610d93565b80601f01602080910402602001604051908101604052809291908181526020018280546102ad90610d93565b80156102fa5780601f106102cf576101008083540402835291602001916102fa565b820191906000526020600020905b8154815290600101906020018083116102dd57829003601f168201915b505050505081565b600061030f33848461055d565b6040518281526001600160a01b0384169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259060200160405180910390a350600192915050565b63ffffffff813b1615155b919050565b6001600160a01b0383166000908152600560209081526040808320338452909152812054828110156103eb5760405162461bcd60e51b815260206004820152602160248201527f7472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636044820152606560f81b60648201526084015b60405180910390fd5b6103ff85336103fa8685610d7c565b61055d565b61040a8585856105df565b506001949350505050565b6006546001600160a01b031633146104695760405162461bcd60e51b815260206004820152601760248201527631b0b63632b91034b9903737ba103a34329037bbb732b960491b60448201526064016103e2565b610473828261078e565b5050565b6001805461028190610d93565b6006546001600160a01b031633146104d85760405162461bcd60e51b815260206004820152601760248201527631b0b63632b91034b9903737ba103a34329037bbb732b960491b60448201526064016103e2565b610473828261086d565b60006104ed33610358565b1561053a5760405162461bcd60e51b815260206004820152601960248201527f63616c6c65722063616e6e6f7420626520636f6e74726163740000000000000060448201526064016103e2565b61040a33868686866109af565b60006105543384846105df565b50600192915050565b6001600160a01b0383166105b35760405162461bcd60e51b815260206004820152601d60248201527f617070726f76652066726f6d20746865207a65726f206164647265737300000060448201526064016103e2565b6001600160a01b0392831660009081526005602090815260408083209490951682529290925291902055565b6001600160a01b0383166106355760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f2061646472657373000060448201526064016103e2565b6001600160a01b03821661068b5760405162461bcd60e51b815260206004820152601c60248201527f7472616e7366657220746f20746865207a65726f20616464726573730000000060448201526064016103e2565b6001600160a01b038316600090815260046020526040902054818110156106f45760405162461bcd60e51b815260206004820152601f60248201527f7472616e7366657220616d6f756e7420657863656564732062616c616e63650060448201526064016103e2565b6106fe8282610d7c565b6001600160a01b038086166000908152600460205260408082209390935590851681529081208054849290610734908490610d64565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161078091815260200190565b60405180910390a350505050565b6001600160a01b0382166107e45760405162461bcd60e51b815260206004820152601860248201527f6d696e7420746f20746865207a65726f2061646472657373000000000000000060448201526064016103e2565b80600360008282546107f69190610d64565b90915550506001600160a01b03821660009081526004602052604081208054839290610823908490610d64565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b6001600160a01b0382166108c35760405162461bcd60e51b815260206004820152601a60248201527f6275726e2066726f6d20746865207a65726f206164647265737300000000000060448201526064016103e2565b6001600160a01b0382166000908152600460205260409020548181101561092c5760405162461bcd60e51b815260206004820152601b60248201527f6275726e20616d6f756e7420657863656564732062616c616e6365000000000060448201526064016103e2565b6109368282610d7c565b6001600160a01b03841660009081526004602052604081209190915560038054849290610964908490610d7c565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b6001600160a01b038516610a055760405162461bcd60e51b815260206004820152601a60248201527f7472616e736665722066726f6d207a65726f206164647265737300000000000060448201526064016103e2565b6000845111610a4e5760405162461bcd60e51b81526020600482015260156024820152747472616e7366657220746f2074686520656d70747960581b60448201526064016103e2565b6000815111610a8e5760405162461bcd60e51b815260206004820152600c60248201526b195b5c1d1e481d185c99d95d60a21b60448201526064016103e2565b610ab385610aa46006546001600160a01b031690565b610aae8587610d64565b6105df565b846001600160a01b03167f922dc141ed1042641d7f53d00c9b504ee06b5007b78fb5edddaeb56d0847422385858585604051610af29493929190610d28565b60405180910390a25050505050565b80356001600160a01b038116811461036357600080fd5b600082601f830112610b28578081fd5b813567ffffffffffffffff80821115610b4357610b43610de4565b604051601f8301601f19908116603f01168101908282118183101715610b6b57610b6b610de4565b81604052838152866020858801011115610b83578485fd5b8360208701602083013792830160200193909352509392505050565b600060208284031215610bb0578081fd5b610bb982610b01565b9392505050565b60008060408385031215610bd2578081fd5b610bdb83610b01565b9150610be960208401610b01565b90509250929050565b600080600060608486031215610c06578081fd5b610c0f84610b01565b9250610c1d60208501610b01565b9150604084013590509250925092565b60008060408385031215610c3f578182fd5b610c4883610b01565b946020939093013593505050565b60008060008060808587031215610c6b578081fd5b843567ffffffffffffffff80821115610c82578283fd5b610c8e88838901610b18565b955060208701359450604087013593506060870135915080821115610cb1578283fd5b50610cbe87828801610b18565b91505092959194509250565b60008151808452815b81811015610cef57602081850181015186830182015201610cd3565b81811115610d005782602083870101525b50601f01601f19169290920160200192915050565b600060208252610bb96020830184610cca565b600060808252610d3b6080830187610cca565b8560208401528460408401528281036060840152610d598185610cca565b979650505050505050565b60008219821115610d7757610d77610dce565b500190565b600082821015610d8e57610d8e610dce565b500390565b600281046001821680610da757607f821691505b60208210811415610dc857634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052604160045260246000fdfea2646970667358221220620e1c008a766c9c7b88a0926c28ff81dd86ff394d2d2c85454d9fbb2ab2cbcc64736f6c63430008020033")
	cCode, _ := hex.DecodeString("608060405234801561001057600080fd5b50600436106100f55760003560e01c806370a0823111610097578063a2ac8a3711610066578063a2ac8a3714610212578063a9059cbb14610225578063b86d529814610238578063dd62ed3e14610249576100f5565b806370a08231146101ac5780638da5cb5b146101cc57806395d89b41146101f75780639dc29fac146101ff576100f5565b806318160ddd116100d357806318160ddd1461014e57806323b872dd14610165578063313ce5671461017857806340c10f1914610197576100f5565b806306fdde03146100fa578063095ea7b314610118578063162790551461013b575b600080fd5b610102610274565b60405161010f9190610d15565b60405180910390f35b61012b610126366004610c2d565b610302565b604051901515815260200161010f565b61012b610149366004610b9f565b610358565b61015760035481565b60405190815260200161010f565b61012b610173366004610bf2565b610368565b6002546101859060ff1681565b60405160ff909116815260200161010f565b6101aa6101a5366004610c2d565b610415565b005b6101576101ba366004610b9f565b60046020526000908152604090205481565b6006546101df906001600160a01b031681565b6040516001600160a01b03909116815260200161010f565b610102610477565b6101aa61020d366004610c2d565b610484565b61012b610220366004610c56565b6104e2565b61012b610233366004610c2d565b610547565b6006546001600160a01b03166101df565b610157610257366004610bc0565b600560209081526000928352604080842090915290825290205481565b6000805461028190610d93565b80601f01602080910402602001604051908101604052809291908181526020018280546102ad90610d93565b80156102fa5780601f106102cf576101008083540402835291602001916102fa565b820191906000526020600020905b8154815290600101906020018083116102dd57829003601f168201915b505050505081565b600061030f33848461055d565b6040518281526001600160a01b0384169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259060200160405180910390a350600192915050565b63ffffffff813b1615155b919050565b6001600160a01b0383166000908152600560209081526040808320338452909152812054828110156103eb5760405162461bcd60e51b815260206004820152602160248201527f7472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636044820152606560f81b60648201526084015b60405180910390fd5b6103ff85336103fa8685610d7c565b61055d565b61040a8585856105df565b506001949350505050565b6006546001600160a01b031633146104695760405162461bcd60e51b815260206004820152601760248201527631b0b63632b91034b9903737ba103a34329037bbb732b960491b60448201526064016103e2565b610473828261078e565b5050565b6001805461028190610d93565b6006546001600160a01b031633146104d85760405162461bcd60e51b815260206004820152601760248201527631b0b63632b91034b9903737ba103a34329037bbb732b960491b60448201526064016103e2565b610473828261086d565b60006104ed33610358565b1561053a5760405162461bcd60e51b815260206004820152601960248201527f63616c6c65722063616e6e6f7420626520636f6e74726163740000000000000060448201526064016103e2565b61040a33868686866109af565b60006105543384846105df565b50600192915050565b6001600160a01b0383166105b35760405162461bcd60e51b815260206004820152601d60248201527f617070726f76652066726f6d20746865207a65726f206164647265737300000060448201526064016103e2565b6001600160a01b0392831660009081526005602090815260408083209490951682529290925291902055565b6001600160a01b0383166106355760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f2061646472657373000060448201526064016103e2565b6001600160a01b03821661068b5760405162461bcd60e51b815260206004820152601c60248201527f7472616e7366657220746f20746865207a65726f20616464726573730000000060448201526064016103e2565b6001600160a01b038316600090815260046020526040902054818110156106f45760405162461bcd60e51b815260206004820152601f60248201527f7472616e7366657220616d6f756e7420657863656564732062616c616e63650060448201526064016103e2565b6106fe8282610d7c565b6001600160a01b038086166000908152600460205260408082209390935590851681529081208054849290610734908490610d64565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161078091815260200190565b60405180910390a350505050565b6001600160a01b0382166107e45760405162461bcd60e51b815260206004820152601860248201527f6d696e7420746f20746865207a65726f2061646472657373000000000000000060448201526064016103e2565b80600360008282546107f69190610d64565b90915550506001600160a01b03821660009081526004602052604081208054839290610823908490610d64565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b6001600160a01b0382166108c35760405162461bcd60e51b815260206004820152601a60248201527f6275726e2066726f6d20746865207a65726f206164647265737300000000000060448201526064016103e2565b6001600160a01b0382166000908152600460205260409020548181101561092c5760405162461bcd60e51b815260206004820152601b60248201527f6275726e20616d6f756e7420657863656564732062616c616e6365000000000060448201526064016103e2565b6109368282610d7c565b6001600160a01b03841660009081526004602052604081209190915560038054849290610964908490610d7c565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b6001600160a01b038516610a055760405162461bcd60e51b815260206004820152601a60248201527f7472616e736665722066726f6d207a65726f206164647265737300000000000060448201526064016103e2565b6000845111610a4e5760405162461bcd60e51b81526020600482015260156024820152747472616e7366657220746f2074686520656d70747960581b60448201526064016103e2565b6000815111610a8e5760405162461bcd60e51b815260206004820152600c60248201526b195b5c1d1e481d185c99d95d60a21b60448201526064016103e2565b610ab385610aa46006546001600160a01b031690565b610aae8587610d64565b6105df565b846001600160a01b03167f922dc141ed1042641d7f53d00c9b504ee06b5007b78fb5edddaeb56d0847422385858585604051610af29493929190610d28565b60405180910390a25050505050565b80356001600160a01b038116811461036357600080fd5b600082601f830112610b28578081fd5b813567ffffffffffffffff80821115610b4357610b43610de4565b604051601f8301601f19908116603f01168101908282118183101715610b6b57610b6b610de4565b81604052838152866020858801011115610b83578485fd5b8360208701602083013792830160200193909352509392505050565b600060208284031215610bb0578081fd5b610bb982610b01565b9392505050565b60008060408385031215610bd2578081fd5b610bdb83610b01565b9150610be960208401610b01565b90509250929050565b600080600060608486031215610c06578081fd5b610c0f84610b01565b9250610c1d60208501610b01565b9150604084013590509250925092565b60008060408385031215610c3f578182fd5b610c4883610b01565b946020939093013593505050565b60008060008060808587031215610c6b578081fd5b843567ffffffffffffffff80821115610c82578283fd5b610c8e88838901610b18565b955060208701359450604087013593506060870135915080821115610cb1578283fd5b50610cbe87828801610b18565b91505092959194509250565b60008151808452815b81811015610cef57602081850181015186830182015201610cd3565b81811115610d005782602083870101525b50601f01601f19169290920160200192915050565b600060208252610bb96020830184610cca565b600060808252610d3b6080830187610cca565b8560208401528460408401528281036060840152610d598185610cca565b979650505050505050565b60008219821115610d7757610d77610dce565b500190565b600082821015610d8e57610d8e610dce565b500390565b600281046001821680610da757607f821691505b60208210811415610dc857634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052604160045260246000fdfea2646970667358221220620e1c008a766c9c7b88a0926c28ff81dd86ff394d2d2c85454d9fbb2ab2cbcc64736f6c63430008020033")
	addr, err := suite.app.IntrarelayerKeeper.DeployContract(suite.ctx, types.ModuleAddress, abiJson, cBin, "FunctionX", "FX", uint8(18))

	suite.Require().NoError(err)

	account := suite.app.EvmKeeper.GetAccount(suite.ctx, addr)
	suite.Require().NotEmpty(account)

	code := suite.app.EvmKeeper.GetCode(suite.ctx, common.BytesToHash(account.CodeHash))
	suite.Require().Equal(cCode, code)
}

func (suite KeeperTestSuite) TestDeployContractERC1967Proxy() {
	suite.SetupTest()

	cABI := contracts.MustGetABI(suite.ctx.BlockHeight(), contracts.ERC1967ProxyType)
	cBin := contracts.MustGetBin(suite.ctx.BlockHeight(), contracts.ERC1967ProxyType)

	addr, err := suite.app.IntrarelayerKeeper.DeployContract(suite.ctx, types.ModuleAddress,
		cABI, cBin, common.HexToAddress(contracts.FIP20UpgradeCodeAddress), []byte{})
	suite.Require().NoError(err)

	account := suite.app.EvmKeeper.GetAccount(suite.ctx, addr)
	suite.Require().NotEmpty(account)

	code := suite.app.EvmKeeper.GetCode(suite.ctx, common.BytesToHash(account.CodeHash))

	cCode := contracts.MustGetCode(suite.ctx.BlockHeight(), contracts.ERC1967ProxyType)
	suite.Require().Equal(cCode, code)
}

func (suite KeeperTestSuite) TestRegisterFIP20() {
	var (
		contractAddr common.Address
		pair         types.TokenPair
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"intrarelaying is disabled globally",
			func() {
				params := types.DefaultParams()
				params.EnableIntrarelayer = false
				suite.app.IntrarelayerKeeper.SetParams(suite.ctx, params)
			},
			false,
		},
		{
			"token FIP20 already registered",
			func() {
				suite.app.IntrarelayerKeeper.SetFIP20Map(suite.ctx, pair.GetFIP20Contract(), pair.GetID())
			},
			false,
		},
		{
			"denom already registered",
			func() {
				suite.app.IntrarelayerKeeper.SetDenomMap(suite.ctx, pair.Denom, pair.GetID())
			},
			false,
		},
		{
			"meta data already stored",
			func() {
				suite.app.IntrarelayerKeeper.CreateCoinMetadata(suite.ctx, contractAddr)
			},
			false,
		},
		{
			"ok",
			func() {},
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			contractAddr = suite.DeployContract(types.ModuleAddress, fip20Name, fip20Symbol, 18)
			suite.Commit()
			coinName := types.CreateDenom(contractAddr.String())
			pair = types.NewTokenPair(contractAddr, coinName, true, types.OWNER_EXTERNAL)

			tc.malleate()

			_, err := suite.app.IntrarelayerKeeper.RegisterFIP20(suite.ctx, contractAddr)
			metadata := suite.app.BankKeeper.GetDenomMetaData(suite.ctx, coinName)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				// Metadata variables
				suite.Require().True(len(metadata.Base) > 0)
				suite.Require().Equal(coinName, metadata.Base)
				//suite.Require().Equal(coinName, metadata.Name)
				suite.Require().Equal(fip20Symbol, metadata.Display)
				//suite.Require().Equal(fip20Symbol, metadata.Symbol)
				// Denom units
				suite.Require().Equal(len(metadata.DenomUnits), 2)
				suite.Require().Equal(coinName, metadata.DenomUnits[0].Denom)
				suite.Require().Equal(uint32(zeroExponent), metadata.DenomUnits[0].Exponent)
				suite.Require().Equal(fip20Symbol, metadata.DenomUnits[1].Denom)
				// Default exponent at contract creation is 18
				suite.Require().Equal(metadata.DenomUnits[1].Exponent, uint32(defaultExponent))
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

func (suite KeeperTestSuite) TestToggleRelay() {
	var (
		contractAddr common.Address
		id           []byte
		pair         types.TokenPair
	)

	testCases := []struct {
		name         string
		malleate     func()
		expPass      bool
		relayEnabled bool
	}{
		{
			"token not registered",
			func() {
				contractAddr = suite.DeployContract(types.ModuleAddress, fip20Name, fip20Symbol, 18)
				suite.Commit()
				pair = types.NewTokenPair(contractAddr, cosmosTokenName, true, types.OWNER_MODULE)
			},
			false,
			false,
		},
		{
			"token not registered - pair not found",
			func() {
				contractAddr = suite.DeployContract(types.ModuleAddress, fip20Name, fip20Symbol, 18)
				suite.Commit()
				pair = types.NewTokenPair(contractAddr, cosmosTokenName, true, types.OWNER_MODULE)
				suite.app.IntrarelayerKeeper.SetFIP20Map(suite.ctx, common.HexToAddress(pair.Fip20Address), pair.GetID())
			},
			false,
			false,
		},
		{
			"disable relay",
			func() {
				contractAddr = suite.setupRegisterFIP20Pair()
				id = suite.app.IntrarelayerKeeper.GetTokenPairID(suite.ctx, contractAddr.String())
				pair, _ = suite.app.IntrarelayerKeeper.GetTokenPair(suite.ctx, id)
			},
			true,
			false,
		},
		{
			"disable and enable relay",
			func() {
				contractAddr = suite.setupRegisterFIP20Pair()
				id = suite.app.IntrarelayerKeeper.GetTokenPairID(suite.ctx, contractAddr.String())
				pair, _ = suite.app.IntrarelayerKeeper.GetTokenPair(suite.ctx, id)
				pair, _ = suite.app.IntrarelayerKeeper.ToggleRelay(suite.ctx, contractAddr.String())
			},
			true,
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset

			tc.malleate()

			var err error
			pair, err = suite.app.IntrarelayerKeeper.ToggleRelay(suite.ctx, contractAddr.String())
			// Request the pair using the GetPairToken func to make sure that is updated on the db
			pair, _ = suite.app.IntrarelayerKeeper.GetTokenPair(suite.ctx, id)
			if tc.expPass {
				suite.Require().NoError(err, tc.name)
				if tc.relayEnabled {
					suite.Require().True(pair.Enabled)
				} else {
					suite.Require().False(pair.Enabled)
				}
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}
