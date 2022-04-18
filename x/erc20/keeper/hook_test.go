package keeper_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibcclienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/core/exported"
	ibctmtypes "github.com/cosmos/cosmos-sdk/x/ibc/light-clients/07-tendermint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/app/ante"
	"github.com/functionx/fx-core/crypto/ethsecp256k1"
	"github.com/functionx/fx-core/server/config"
	"github.com/functionx/fx-core/tests"
	fxtypes "github.com/functionx/fx-core/types"
	bsctypes "github.com/functionx/fx-core/x/bsc/types"
	"github.com/functionx/fx-core/x/crosschain"
	crosschainkeeper "github.com/functionx/fx-core/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
	"github.com/functionx/fx-core/x/erc20/keeper"
	"github.com/functionx/fx-core/x/erc20/types"
	evmkeeper "github.com/functionx/fx-core/x/evm/keeper"
	evm "github.com/functionx/fx-core/x/evm/types"
	"github.com/functionx/fx-core/x/gravity"
	gravitytypes "github.com/functionx/fx-core/x/gravity/types"
	ibctransfertypes "github.com/functionx/fx-core/x/ibc/applications/transfer/types"
)

type IBCTransferSimulate struct {
	T *testing.T
}

func (it *IBCTransferSimulate) SendTransfer(ctx sdk.Context, sourcePort, sourceChannel string, token sdk.Coin, sender sdk.AccAddress,
	receiver string, timeoutHeight ibcclienttypes.Height, timeoutTimestamp uint64, router string, fee sdk.Coin) error {
	return nil
}

func (it *IBCTransferSimulate) Transfer(goCtx context.Context, msg *ibctransfertypes.MsgTransfer) (*ibctransfertypes.MsgTransferResponse, error) {
	return &ibctransfertypes.MsgTransferResponse{}, nil
}

func (it *IBCTransferSimulate) GetRouter() *ibctransfertypes.Router {
	router := ibctransfertypes.NewRouter()

	return router
}

type IBCChannelSimulate struct {
}

func (ic *IBCChannelSimulate) GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, exported.ClientState, error) {
	return "", &ibctmtypes.ClientState{
		ChainId:         "fxcore",
		TrustLevel:      ibctmtypes.Fraction{},
		TrustingPeriod:  0,
		UnbondingPeriod: 0,
		MaxClockDrift:   0,
		FrozenHeight: ibcclienttypes.Height{
			RevisionHeight: 1000,
			RevisionNumber: 1000,
		},
		LatestHeight: ibcclienttypes.Height{
			RevisionHeight: 10,
			RevisionNumber: 10,
		},
		ProofSpecs:                   nil,
		UpgradePath:                  nil,
		AllowUpdateAfterExpiry:       false,
		AllowUpdateAfterMisbehaviour: false,
	}, nil
}

func (ic *IBCChannelSimulate) GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	return 1, true
}

var (
	wfxMetadata = banktypes.Metadata{
		Description: "Wrap Function X",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "FX",
				Exponent: 0,
				Aliases:  nil,
			},
			{
				Denom:    "WFX",
				Exponent: 18,
				Aliases:  nil,
			},
		},
		Base:    "FX",
		Display: "WFX",
	}

	purseMetadata = banktypes.Metadata{
		Description: "Pundi X Purse Token",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    purseDenom,
				Exponent: 0,
				Aliases:  nil,
			},
			{
				Denom:    "PURSE",
				Exponent: 18,
				Aliases:  nil,
			},
		},
		Base:    purseDenom,
		Display: "PURSE",
	}

	purseDenom = "ibc/5BAB702195E8411500FC0256EBC717AB9C7DC039D7FFC6E1A471022AA939E600"
)

func TestHookChainGravity(t *testing.T) {
	fxcore, validators, _, delegateAddressArr := initTest(t)
	ctx := fxcore.BaseApp.NewContext(false, tmproto.Header{ProposerAddress: validators[0].Address, Height: fxtypes.EvmSupportBlock()})
	require.NoError(t, InitEvmModuleParams(ctx, &fxcore.Erc20Keeper, true))

	pair, err := fxcore.Erc20Keeper.RegisterCoin(ctx, wfxMetadata)
	require.NoError(t, err)

	val := validators[0]
	validator := GetValidator(t, fxcore, val)[0]
	del := delegateAddressArr[0]

	ctx = ctx.WithBlockHeight(504000)

	signer1, addr1 := privateSigner()
	_, addr2 := privateSigner()
	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100))
	err = fxcore.BankKeeper.SendCoins(ctx, del, sdk.AccAddress(addr1.Bytes()), sdk.NewCoins(sdk.NewCoin(fxtypes.MintDenom, amt)))
	require.NoError(t, err)

	ctx = testInitGravity(t, ctx, fxcore, validator.GetOperator(), addr1.Bytes(), addr2)

	balances := fxcore.BankKeeper.GetAllBalances(ctx, addr1.Bytes())
	_ = balances

	err = fxcore.Erc20Keeper.RelayConvertCoin(ctx, addr1.Bytes(), addr1, sdk.NewCoin(fxtypes.MintDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	require.NoError(t, err)

	balanceOf, err := fxcore.Erc20Keeper.BalanceOf(ctx, pair.GetERC20Contract(), addr1)
	require.NoError(t, err)
	_ = balanceOf

	token := pair.GetERC20Contract()
	crossChainTarget := fmt.Sprintf("%s%s", fxtypes.FIP20TransferToChainPrefix, gravitytypes.ModuleName)
	transferChainData := packTransferCrossData(t, ctx, fxcore.Erc20Keeper, addr2.String(), big.NewInt(1e18), big.NewInt(1e18), crossChainTarget)
	sendEthTx(t, ctx, fxcore, signer1, addr1, token, transferChainData)

	transactions := fxcore.GravityKeeper.GetPoolTransactions(ctx)
	_ = transactions
}

func TestHookChainBSC(t *testing.T) {
	fxcore, validators, genesisAccount, delegateAddressArr := initTest(t)
	ctx := fxcore.BaseApp.NewContext(false, tmproto.Header{ProposerAddress: validators[0].Address, Height: fxtypes.EvmSupportBlock()})
	// 1. init evm + erc20 module params
	require.NoError(t, InitEvmModuleParams(ctx, &fxcore.Erc20Keeper, true))

	// 2. register erc20 module coin
	pair, err := fxcore.Erc20Keeper.RegisterCoin(ctx, purseMetadata)
	require.NoError(t, err)

	purseID := fxcore.Erc20Keeper.GetDenomMap(ctx, purseMetadata.Base)
	require.NotEmpty(t, purseID)
	purseTokenPair, found := fxcore.Erc20Keeper.GetTokenPair(ctx, purseID)
	require.True(t, found)
	require.NotNil(t, purseTokenPair)
	require.NotEmpty(t, purseTokenPair.GetErc20Address())

	require.Equal(t, types.TokenPair{
		Erc20Address:  purseTokenPair.GetErc20Address(),
		Denom:         purseMetadata.GetBase(),
		Enabled:       true,
		ContractOwner: types.OWNER_MODULE,
	}, purseTokenPair)

	fip20, err := fxcore.Erc20Keeper.QueryERC20(ctx, pair.GetERC20Contract())
	require.NoError(t, err)
	_ = fip20
	//t.Log("fip20", fip20.Name, fip20.Symbol, fip20.Decimals)

	del := delegateAddressArr[0]
	ga := genesisAccount[0]

	ctx = ctx.WithBlockHeight(504000)

	signer1, addr1 := privateSigner()
	_, addr2 := privateSigner()
	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100))

	err = fxcore.BankKeeper.SendCoins(ctx, ga.GetAddress(), addr1.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.MintDenom, amt), sdk.NewCoin(purseDenom, amt)))
	require.NoError(t, err)

	ctx = testInitBscCrossChain(t, ctx, fxcore, del, addr1.Bytes(), addr2)

	balances := fxcore.BankKeeper.GetAllBalances(ctx, addr1.Bytes())
	_ = balances

	err = fxcore.Erc20Keeper.RelayConvertCoin(ctx, addr1.Bytes(), addr1, sdk.NewCoin(purseDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	require.NoError(t, err)

	balanceOf, err := fxcore.Erc20Keeper.BalanceOf(ctx, pair.GetERC20Contract(), addr1)
	require.NoError(t, err)
	_ = balanceOf

	token := pair.GetERC20Contract()
	crossChainTarget := fmt.Sprintf("%s%s", fxtypes.FIP20TransferToChainPrefix, bsctypes.ModuleName)
	transferChainData := packTransferCrossData(t, ctx, fxcore.Erc20Keeper, addr2.String(), big.NewInt(1e18), big.NewInt(1e18), crossChainTarget)
	sendEthTx(t, ctx, fxcore, signer1, addr1, token, transferChainData)

	transactions := fxcore.BscKeeper.GetUnbatchedTransactions(ctx)
	_ = transactions
}

func TestHookIBC(t *testing.T) {
	fxcore, validators, _, delegateAddressArr := initTest(t)
	ctx := fxcore.BaseApp.NewContext(false, tmproto.Header{ProposerAddress: validators[0].Address, Height: fxtypes.EvmSupportBlock()})
	require.NoError(t, InitEvmModuleParams(ctx, &fxcore.Erc20Keeper, true))

	pair, err := fxcore.Erc20Keeper.RegisterCoin(ctx, wfxMetadata)
	require.NoError(t, err)

	//validator := GetValidator(t, app, val)[0]
	//val := validators[0]
	del := delegateAddressArr[0]

	ctx = ctx.WithBlockHeight(504000)

	signer1, addr1 := privateSigner()
	//_, addr2 := privateSigner()
	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100))
	err = fxcore.BankKeeper.SendCoins(ctx, del, sdk.AccAddress(addr1.Bytes()), sdk.NewCoins(sdk.NewCoin(fxtypes.MintDenom, amt)))
	require.NoError(t, err)

	balances := fxcore.BankKeeper.GetAllBalances(ctx, addr1.Bytes())
	_ = balances

	err = fxcore.Erc20Keeper.RelayConvertCoin(ctx, addr1.Bytes(), addr1, sdk.NewCoin(fxtypes.MintDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	require.NoError(t, err)

	balanceOf, err := fxcore.Erc20Keeper.BalanceOf(ctx, pair.GetERC20Contract(), addr1)
	require.NoError(t, err)
	_ = balanceOf

	//reset ibc
	fxcore.Erc20Keeper.SetIBCTransferKeeperForTest(&IBCTransferSimulate{T: t})
	fxcore.Erc20Keeper.SetIBCChannelKeeperForTest(&IBCChannelSimulate{})

	evmHooks := evmkeeper.NewMultiEvmHooks(fxcore.Erc20Keeper.Hooks())
	fxcore.EvmKeeper = fxcore.EvmKeeper.SetHooksForTest(evmHooks)

	token := pair.GetERC20Contract()
	ibcTarget := fmt.Sprintf("%s%s", fxtypes.FIP20TransferToChainPrefix, "px/transfer/channel-0")
	transferIBCData := packTransferCrossData(t, ctx, fxcore.Erc20Keeper, "px16u6kjunrcxkvaln9aetxwjpruply3sgwpr9z8u", big.NewInt(1e18), big.NewInt(0), ibcTarget)
	sendEthTx(t, ctx, fxcore, signer1, addr1, token, transferIBCData)
}

func packTransferCrossData(t *testing.T, ctx sdk.Context, k keeper.Keeper, to string, amount, fee *big.Int, target string) []byte {
	fip20 := fxtypes.GetERC20(ctx.BlockHeight())
	var targetIBCByte32 [32]byte
	copy(targetIBCByte32[:], target)
	pack, err := fip20.ABI.Pack("transferCrossChain", to, amount, fee, targetIBCByte32)
	require.NoError(t, err)
	return pack
}

func sendEthTx(t *testing.T, ctx sdk.Context, fxcore *app.App,
	signer keyring.Signer, from, contract common.Address, data []byte) {

	chainID := fxcore.EvmKeeper.ChainID()

	args, err := json.Marshal(&evm.TransactionArgs{To: &contract, From: &from, Data: (*hexutil.Bytes)(&data)})
	require.NoError(t, err)
	res, err := fxcore.EvmKeeper.EstimateGas(sdk.WrapSDKContext(ctx), &evm.EthCallRequest{
		Args:   args,
		GasCap: uint64(config.DefaultGasCap),
	})
	require.NoError(t, err)

	nonce, err := fxcore.AccountKeeper.GetSequence(ctx, from.Bytes())
	require.NoError(t, err)

	ercTransferTx := evm.NewTx(
		chainID,
		nonce,
		&contract,
		nil,
		res.Gas,
		nil,
		fxcore.FeeMarketKeeper.GetBaseFee(ctx),
		big.NewInt(1),
		data,
		&ethtypes.AccessList{}, // accesses
	)

	ercTransferTx.From = from.String()
	err = ercTransferTx.Sign(ethtypes.LatestSignerForChainID(chainID), signer)
	require.NoError(t, err)

	options := ante.HandlerOptions{
		AccountKeeper:   fxcore.AccountKeeper,
		BankKeeper:      fxcore.BankKeeper,
		EvmKeeper:       fxcore.EvmKeeper,
		SignModeHandler: app.MakeEncodingConfig().TxConfig.SignModeHandler(),
		SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
	}
	require.NoError(t, options.Validate())
	handler := ante.NewAnteHandler(options)

	clientCtx := client.Context{}.WithTxConfig(app.MakeEncodingConfig().TxConfig)
	params := fxcore.EvmKeeper.GetParams(ctx)
	tx, err := ercTransferTx.BuildTx(clientCtx.TxConfig.NewTxBuilder(), params.EvmDenom)
	require.NoError(t, err)
	ctx, err = handler(ctx, tx, false)
	require.NoError(t, err)

	rsp, err := fxcore.EvmKeeper.EthereumTx(sdk.WrapSDKContext(ctx), ercTransferTx)
	require.NoError(t, err)
	require.Empty(t, rsp.VmError)
}

func privateSigner() (keyring.Signer, common.Address) {
	// account key
	priKey := NewPriKey()
	//ethsecp256k1.GenerateKey()
	ethPriv := &ethsecp256k1.PrivKey{Key: priKey.Bytes()}

	return tests.NewSigner(ethPriv), common.BytesToAddress(ethPriv.PubKey().Address())
}

func initTest(t *testing.T) (*app.App, []*tmtypes.Validator, authtypes.GenesisAccounts, []sdk.AccAddress) {
	initBalances := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(20000))
	validator, genesisAccounts, balances := app.GenerateGenesisValidator(1,
		sdk.NewCoins(sdk.NewCoin(fxtypes.MintDenom, initBalances),
			sdk.NewCoin("eth0xd9EEd31F5731DfC3Ca18f09B487e200F50a6343B", initBalances),
			sdk.NewCoin("ibc/4757BC3AA2C696F7083C825BD3951AE3D1631F2A272EA7AFB9B3E1CCCA8560D4", initBalances),
			sdk.NewCoin(purseDenom, initBalances)))

	fxcore := app.SetupWithGenesisValSet(t, validator, genesisAccounts, balances...)
	ctx := fxcore.BaseApp.NewContext(false, tmproto.Header{})
	delegateAddressArr := app.AddTestAddrsIncremental(fxcore, ctx, 1, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(1000000000)))

	fxtypes.ChangeNetworkForTest(fxtypes.NetworkDevnet())
	return fxcore, validator.Validators, genesisAccounts, delegateAddressArr
}

func GetValidator(t *testing.T, fxcore *app.App, vals ...*tmtypes.Validator) []stakingtypes.Validator {
	ctx := fxcore.BaseApp.NewContext(false, tmproto.Header{})
	validators := make([]stakingtypes.Validator, 0, len(vals))
	for _, val := range vals {
		validator, found := fxcore.StakingKeeper.GetValidator(ctx, val.Address.Bytes())
		require.True(t, found)
		validators = append(validators, validator)
	}
	return validators
}

var (
	FxOriginatedTokenContract = common.HexToAddress("0x0000000000000000000000000000000000000000")
	BSCBridgeTokenContract    = common.HexToAddress("0x0000000000000000000000000000000000000001")
)

func testInitGravity(t *testing.T, ctx sdk.Context, fxcore *app.App, val sdk.ValAddress, orch sdk.AccAddress, addr common.Address) sdk.Context {
	fxcore.GravityKeeper.SetOrchestratorValidator(ctx, val, orch)
	fxcore.GravityKeeper.SetEthAddressForValidator(ctx, val, addr.String())

	testValSetUpdateClaim(t, ctx, fxcore, orch, addr)

	testFxOriginatedTokenClaim(t, ctx, fxcore, orch)

	gravity.EndBlocker(ctx, fxcore.GravityKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	return ctx
}

func testFxOriginatedTokenClaim(t *testing.T, ctx sdk.Context, fxcore *app.App, orch sdk.AccAddress) {
	msg := &gravitytypes.MsgFxOriginatedTokenClaim{
		EventNonce:    2,
		BlockHeight:   uint64(ctx.BlockHeight()),
		TokenContract: FxOriginatedTokenContract.String(),
		Name:          "Function X",
		Symbol:        "FX",
		Decimals:      18,
		Orchestrator:  orch.String(),
	}

	any, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = fxcore.GravityKeeper.Attest(ctx, msg, any)
	require.NoError(t, err)
}

func testValSetUpdateClaim(t *testing.T, ctx sdk.Context, fxcore *app.App, orch sdk.AccAddress, addr common.Address) {
	msg := &gravitytypes.MsgValsetUpdatedClaim{
		EventNonce:  1,
		BlockHeight: uint64(ctx.BlockHeight()),
		ValsetNonce: 0,
		Members: []*gravitytypes.BridgeValidator{
			{
				Power:      uint64(math.MaxUint32),
				EthAddress: addr.String(),
			},
		},
		Orchestrator: orch.String(),
	}

	for _, member := range msg.Members {
		memberVal := fxcore.GravityKeeper.GetValidatorByEthAddress(ctx, member.EthAddress)
		require.NotEmpty(t, memberVal)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = fxcore.GravityKeeper.Attest(ctx, msg, any)
	require.NoError(t, err)
}

func testInitBscCrossChain(t *testing.T, ctx sdk.Context, fxcore *app.App, oracleAddress, orchestratorAddr sdk.AccAddress, externalAddress common.Address) sdk.Context {
	deposit := sdk.NewCoin(fxtypes.MintDenom, sdk.NewIntFromBigInt(big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18))))
	err := fxcore.BankKeeper.SendCoinsFromAccountToModule(ctx, oracleAddress, fxcore.BscKeeper.GetModuleName(), sdk.NewCoins(deposit))
	require.NoError(t, err)

	testBSCParamsProposal(t, ctx, fxcore, oracleAddress)

	oracle := crosschaintypes.Oracle{
		OracleAddress:       oracleAddress.String(),
		OrchestratorAddress: orchestratorAddr.String(),
		ExternalAddress:     externalAddress.String(),
		DepositAmount:       deposit,
		StartHeight:         ctx.BlockHeight(),
		Jailed:              false,
		JailedHeight:        0,
	}
	// save oracle
	fxcore.BscKeeper.SetOracle(ctx, oracle)

	fxcore.BscKeeper.SetOracleByOrchestrator(ctx, oracleAddress, orchestratorAddr)
	// set the ethereum address
	fxcore.BscKeeper.SetExternalAddressForOracle(ctx, oracleAddress, externalAddress.String())
	// save total deposit amount
	totalDeposit := fxcore.BscKeeper.GetTotalDeposit(ctx)
	fxcore.BscKeeper.SetTotalDeposit(ctx, totalDeposit.Add(deposit))

	fxcore.BscKeeper.CommonSetOracleTotalPower(ctx)

	testBSCOracleSetUpdateClaim(t, ctx, fxcore, orchestratorAddr, externalAddress)

	testBSCBridgeTokenClaim(t, ctx, fxcore, orchestratorAddr)

	crosschain.EndBlocker(ctx, fxcore.BscKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	return ctx
}

func testBSCParamsProposal(t *testing.T, ctx sdk.Context, fxcore *app.App, oracles sdk.AccAddress) {
	slashFraction, _ := sdk.NewDecFromStr("0.001")
	oracleSetUpdatePowerChangePercent, _ := sdk.NewDecFromStr("0.1")
	proposal := &crosschaintypes.InitCrossChainParamsProposal{
		Title:       "bsc cross chain",
		Description: "bsc cross chain init",
		Params: &crosschaintypes.Params{
			GravityId:                         "fx-bsc-bridge",
			SignedWindow:                      20000,
			ExternalBatchTimeout:              86400000,
			AverageBlockTime:                  1000,
			AverageExternalBlockTime:          3000,
			SlashFraction:                     slashFraction,
			OracleSetUpdatePowerChangePercent: oracleSetUpdatePowerChangePercent,
			IbcTransferTimeoutHeight:          20000,
			Oracles:                           []string{oracles.String()},
			DepositThreshold:                  sdk.NewCoin(fxtypes.MintDenom, sdk.NewIntFromBigInt(big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18)))),
		},
		ChainName: fxcore.BscKeeper.GetModuleName(),
	}

	k := crosschainkeeper.EthereumMsgServer{Keeper: fxcore.BscKeeper}
	err := k.HandleInitCrossChainParamsProposal(ctx, proposal)
	require.NoError(t, err)
}

func testBSCBridgeTokenClaim(t *testing.T, ctx sdk.Context, fxcore *app.App, orchAddr sdk.AccAddress) {
	msg := &crosschaintypes.MsgBridgeTokenClaim{
		EventNonce:    2,
		BlockHeight:   uint64(ctx.BlockHeight()),
		TokenContract: BSCBridgeTokenContract.String(),
		Name:          "PURSE Token",
		Symbol:        "PURSE",
		Decimals:      18,
		Orchestrator:  orchAddr.String(),
		ChannelIbc:    hex.EncodeToString([]byte("transfer/channel-0")),
		ChainName:     fxcore.BscKeeper.GetModuleName(),
	}

	any, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = fxcore.BscKeeper.Attest(ctx, msg, any)
	require.NoError(t, err)
}

func testBSCOracleSetUpdateClaim(t *testing.T, ctx sdk.Context, fxcore *app.App, orch sdk.AccAddress, addr common.Address) {
	msg := &crosschaintypes.MsgOracleSetUpdatedClaim{
		EventNonce:     1,
		BlockHeight:    uint64(ctx.BlockHeight()),
		OracleSetNonce: 0,
		Members: crosschaintypes.BridgeValidators{
			{
				Power:           uint64(math.MaxUint32),
				ExternalAddress: addr.String(),
			},
		},
		Orchestrator: orch.String(),
		ChainName:    fxcore.BscKeeper.GetModuleName(),
	}
	for _, member := range msg.Members {
		_, found := fxcore.BscKeeper.GetOracleByExternalAddress(ctx, member.ExternalAddress)
		require.True(t, found)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = fxcore.BscKeeper.Attest(ctx, msg, any)
	require.NoError(t, err)
}