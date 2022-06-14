package crosschain_test

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/functionx/fx-core/x/crosschain/keeper"

	"github.com/functionx/fx-core/app/helpers"

	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/functionx/fx-core/app"
	"github.com/functionx/fx-core/x/crosschain"
	"github.com/functionx/fx-core/x/crosschain/types"
)

type MsgHandlerTestSuite struct {
	suite.Suite
	sync.Mutex

	app            *app.App
	ctx            sdk.Context
	oracles        []sdk.AccAddress
	bridgers       []sdk.AccAddress
	externals      []*ecdsa.PrivateKey
	validator      []sdk.ValAddress
	chainName      string
	delegateAmount sdk.Int
}

func TestMsgHandlerTestSuite(t *testing.T) {
	methodFinder := reflect.TypeOf(new(MsgHandlerTestSuite))
	for i := 0; i < methodFinder.NumMethod(); i++ {
		method := methodFinder.Method(i)
		if !strings.HasPrefix(method.Name, "Test") {
			continue
		}
		t.Run(method.Name, func(subT *testing.T) {
			mySuite := new(MsgHandlerTestSuite)
			mySuite.SetT(t)
			mySuite.SetupTest()
			method.Func.Call([]reflect.Value{reflect.ValueOf(mySuite)})
		})
	}
}

func (suite *MsgHandlerTestSuite) Handler() sdk.Handler {
	return crosschain.NewHandler(suite.app.CrosschainKeeper)
}

func (suite *MsgHandlerTestSuite) Keeper() keeper.Keeper {
	return suite.app.BscKeeper
}

func (suite *MsgHandlerTestSuite) SetupTest() {
	valSet, valAccounts, valBalances := helpers.GenerateGenesisValidator(types.MaxOracleSize, sdk.Coins{})
	suite.app = helpers.SetupWithGenesisValSet(suite.T(), valSet, valAccounts, valBalances...)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
	suite.oracles = helpers.AddTestAddrs(suite.app, suite.ctx, types.MaxOracleSize, sdk.NewInt(300*1e3).MulRaw(1e18))
	suite.bridgers = helpers.AddTestAddrs(suite.app, suite.ctx, types.MaxOracleSize, sdk.NewInt(300*1e3).MulRaw(1e18))
	suite.externals = genEthKey(types.MaxOracleSize)
	suite.delegateAmount = sdk.NewInt(10 * 1e3).MulRaw(1e18)
	for i := 0; i < types.MaxOracleSize; i++ {
		suite.validator = append(suite.validator, valAccounts[i].GetAddress().Bytes())
	}
	suite.chainName = "bsc"

	proposalOracle := &types.ProposalOracle{}
	for _, oracle := range suite.oracles {
		proposalOracle.Oracles = append(proposalOracle.Oracles, oracle.String())
	}
	suite.Keeper().SetProposalOracle(suite.ctx, proposalOracle)
}

// 1. Test MsgCreateOracleBridger
func (suite *MsgHandlerTestSuite) TestMsgCreateOracleBridger() {

	// 1. sender not in chain oracle
	notOracleMsg := &types.MsgCreateOracleBridger{
		OracleAddress:    suite.bridgers[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(100000)},
		ChainName:        suite.chainName,
	}
	_, err := suite.Handler()(suite.ctx, notOracleMsg)
	require.ErrorIs(suite.T(), err, types.ErrNoFoundOracle)

	// 2. stake denom not match chain params stake denom
	notMatchStakeDenomMsg := &types.MsgCreateOracleBridger{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount:   sdk.Coin{Denom: "abctoken", Amount: sdk.NewInt(100000)},
		ChainName:        suite.chainName,
	}
	_, err = suite.Handler()(suite.ctx, notMatchStakeDenomMsg)
	require.ErrorIs(suite.T(), err, types.ErrInvalid)
	require.EqualValues(suite.T(), fmt.Sprintf(
		"delegate denom got %s, expected %s: %s",
		notMatchStakeDenomMsg.DelegateAmount.Denom, fxtypes.DefaultDenom, types.ErrInvalid), err.Error())

	// 3. insufficient stake amount msg.
	belowMinimumStakeAmountMsg := &types.MsgCreateOracleBridger{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(100000)},
		ChainName:        suite.chainName,
	}
	_, err = suite.Handler()(suite.ctx, belowMinimumStakeAmountMsg)
	require.ErrorIs(suite.T(), types.ErrDelegateAmountBelowMinimum, err)
	require.EqualValues(suite.T(), types.ErrDelegateAmountBelowMinimum.Error(), err.Error())

	// 4. success msg
	normalMsg := &types.MsgCreateOracleBridger{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[1].String(),
		DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: suite.delegateAmount},
		ChainName:        suite.chainName,
	}
	_, err = suite.Handler()(suite.ctx, normalMsg)
	require.NoError(suite.T(), err)

	// 5. oracle duplicate set bridger
	oracleDuplicateSetOrchestratorMsg := &types.MsgCreateOracleBridger{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(100000)},
		ChainName:        suite.chainName,
	}
	_, err = suite.Handler()(suite.ctx, oracleDuplicateSetOrchestratorMsg)
	require.ErrorIs(suite.T(), types.ErrInvalid, err)
	require.EqualValues(suite.T(), fmt.Sprintf("oracle existed bridger address: %s", types.ErrInvalid.Error()), err.Error())

	// 6. Set the same bridger address for different Oracle databases
	duplicateSetOrchestratorMsg := &types.MsgCreateOracleBridger{
		OracleAddress:    suite.oracles[1].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(100000)},
		ChainName:        suite.chainName,
	}
	_, err = suite.Handler()(suite.ctx, duplicateSetOrchestratorMsg)
	require.ErrorIs(suite.T(), types.ErrInvalid, err)
	require.EqualValues(suite.T(), fmt.Sprintf("bridger address is bound to oracle: %s", types.ErrInvalid.Error()), err.Error())

	// 7. Set the same external address for different Oracle databases
	duplicateSetExternalAddressMsg := &types.MsgCreateOracleBridger{
		OracleAddress:    suite.oracles[1].String(),
		BridgerAddress:   suite.bridgers[1].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(100000)},
		ChainName:        suite.chainName,
	}
	_, err = suite.Handler()(suite.ctx, duplicateSetExternalAddressMsg)
	require.ErrorIs(suite.T(), types.ErrInvalid, err)
	require.EqualValues(suite.T(), fmt.Sprintf("external address is bound to oracle: %s", types.ErrInvalid.Error()), err.Error())

	// 8. Margin is not allowed to be submitted more than ten times the threshold
	depositAmountBelowMaximumMsg := &types.MsgCreateOracleBridger{
		OracleAddress:    suite.oracles[1].String(),
		BridgerAddress:   suite.bridgers[1].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[1].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: suite.delegateAmount.Mul(sdk.NewInt(10).Add(sdk.NewInt(1)))},
		ChainName:        suite.chainName,
	}
	_, err = suite.Handler()(suite.ctx, depositAmountBelowMaximumMsg)
	require.ErrorIs(suite.T(), types.ErrDelegateAmountBelowMaximum, err)
	require.EqualValues(suite.T(), types.ErrDelegateAmountBelowMaximum.Error(), err.Error())

	// 9. success msg
	normalMsgOracle2 := &types.MsgCreateOracleBridger{
		OracleAddress:    suite.oracles[1].String(),
		BridgerAddress:   suite.bridgers[1].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[1].PublicKey).Hex(),
		ValidatorAddress: suite.validator[1].String(),
		DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: suite.delegateAmount},
		ChainName:        suite.chainName,
	}
	_, err = suite.Handler()(suite.ctx, normalMsgOracle2)
	require.NoError(suite.T(), err)
}

// 2. Test MsgAddOracleDelegate
func (suite *MsgHandlerTestSuite) TestMsgAddOracleDelegate() {

	normalMsg := &types.MsgCreateOracleBridger{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: suite.delegateAmount},
		ChainName:        suite.chainName,
	}
	_, err := suite.Handler()(suite.ctx, normalMsg)
	require.NoError(suite.T(), err)

	denomNotMatchMsg := &types.MsgAddOracleDelegate{
		OracleAddress: suite.oracles[0].String(),
		Amount: sdk.Coin{
			Denom:  "abc",
			Amount: suite.delegateAmount,
		},
		ChainName: suite.chainName,
	}
	_, err = suite.Handler()(suite.ctx, denomNotMatchMsg)
	require.ErrorIs(suite.T(), err, types.ErrInvalid)
	require.EqualValues(suite.T(), fmt.Sprintf("delegate denom got %s, expected %s: %s", denomNotMatchMsg.Amount.Denom, fxtypes.DefaultDenom, types.ErrInvalid), err.Error())

	notOracleMsg := &types.MsgAddOracleDelegate{
		OracleAddress: suite.bridgers[0].String(),
		Amount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: suite.delegateAmount,
		},
		ChainName: suite.chainName,
	}
	_, err = suite.Handler()(suite.ctx, notOracleMsg)
	require.ErrorIs(suite.T(), types.ErrNoFoundOracle, err)
	require.EqualValues(suite.T(), types.ErrNoFoundOracle.Error(), err.Error())

	notSetBridgerOracleMsg := &types.MsgAddOracleDelegate{
		OracleAddress: suite.oracles[2].String(),
		Amount: sdk.Coin{
			Denom:  fxtypes.DefaultDenom,
			Amount: suite.delegateAmount,
		},
		ChainName: suite.chainName,
	}
	_, err = suite.Handler()(suite.ctx, notSetBridgerOracleMsg)
	require.ErrorIs(suite.T(), types.ErrNoFoundOracle, err)
	require.EqualValues(suite.T(), types.ErrNoFoundOracle.Error(), err.Error())

	depositAmountBelowMaximumMsg := &types.MsgAddOracleDelegate{
		OracleAddress: suite.oracles[0].String(),
		Amount:        sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: suite.delegateAmount.Mul(sdk.NewInt(9)).Add(sdk.NewInt(1))},
		ChainName:     suite.chainName,
	}
	_, err = suite.Handler()(suite.ctx, depositAmountBelowMaximumMsg)
	require.ErrorIs(suite.T(), types.ErrDelegateAmountBelowMaximum, err)
	require.EqualValues(suite.T(), types.ErrDelegateAmountBelowMaximum.Error(), err.Error())

	normalAddStakeMsg := &types.MsgAddOracleDelegate{
		OracleAddress: suite.oracles[0].String(),
		Amount:        sdk.NewCoin(fxtypes.DefaultDenom, suite.delegateAmount),
		ChainName:     suite.chainName,
	}
	_, err = suite.Handler()(suite.ctx, normalAddStakeMsg)
	require.NoError(suite.T(), err)
}

func (suite *MsgHandlerTestSuite) TestMsgSetOracleSetConfirm() {

	normalMsg := &types.MsgCreateOracleBridger{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: suite.delegateAmount},
		ChainName:        suite.chainName,
	}
	_, err := suite.Handler()(suite.ctx, normalMsg)
	require.NoError(suite.T(), err)

	latestOracleSetNonce := suite.Keeper().GetLatestOracleSetNonce(suite.ctx)
	require.EqualValues(suite.T(), 0, latestOracleSetNonce)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	latestOracleSetNonce = suite.Keeper().GetLatestOracleSetNonce(suite.ctx)
	require.EqualValues(suite.T(), 1, latestOracleSetNonce)

	require.True(suite.T(), suite.Keeper().HasOracleSetRequest(suite.ctx, 1))

	require.False(suite.T(), suite.Keeper().HasOracleSetRequest(suite.ctx, 2))

	nonce1OracleSet := suite.Keeper().GetOracleSet(suite.ctx, 1)
	require.EqualValues(suite.T(), uint64(1), nonce1OracleSet.Nonce)
	require.EqualValues(suite.T(), uint64(2), nonce1OracleSet.Height)
	require.EqualValues(suite.T(), 1, len(nonce1OracleSet.Members))
	require.EqualValues(suite.T(), crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(), nonce1OracleSet.Members[0].ExternalAddress)
	require.EqualValues(suite.T(), math.MaxUint32, nonce1OracleSet.Members[0].Power)

	var gravityId string
	require.NotPanics(suite.T(), func() {
		gravityId = suite.Keeper().GetGravityID(suite.ctx)
	})
	require.EqualValues(suite.T(), "fx-bridge-bsc", gravityId)
	checkpoint, err := nonce1OracleSet.GetCheckpoint(gravityId)

	external1Signature, err := types.NewEthereumSignature(checkpoint, suite.externals[0])
	require.NoError(suite.T(), err)
	external2Signature, err := types.NewEthereumSignature(checkpoint, suite.externals[1])
	require.NoError(suite.T(), err)
	errMsgData := []struct {
		name      string
		msg       *types.MsgOracleSetConfirm
		err       error
		errReason string
	}{
		{
			name: "Error oracleSet nonce",
			msg: &types.MsgOracleSetConfirm{
				Nonce:           0,
				BridgerAddress:  suite.bridgers[0].String(),
				ExternalAddress: crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
				Signature:       hex.EncodeToString(external1Signature),
				ChainName:       suite.chainName,
			},
			err:       types.ErrInvalid,
			errReason: fmt.Sprintf("couldn't find oracleSet: %s", types.ErrInvalid),
		},
		{
			name: "not oracle msg",
			msg: &types.MsgOracleSetConfirm{
				Nonce:           nonce1OracleSet.Nonce,
				BridgerAddress:  suite.bridgers[0].String(),
				ExternalAddress: crypto.PubkeyToAddress(suite.externals[1].PublicKey).Hex(),
				Signature:       hex.EncodeToString(external1Signature),
				ChainName:       suite.chainName,
			},
			err:       types.ErrNoFoundOracle,
			errReason: fmt.Sprintf("%s", types.ErrNoFoundOracle),
		},
		{
			name: "sign not match external-1  external-sign-2",
			msg: &types.MsgOracleSetConfirm{
				Nonce:           nonce1OracleSet.Nonce,
				BridgerAddress:  suite.bridgers[0].String(),
				ExternalAddress: crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
				Signature:       hex.EncodeToString(external2Signature),
				ChainName:       suite.chainName,
			},
			err:       types.ErrInvalid,
			errReason: fmt.Sprintf("signature verification failed expected sig by %s with checkpoint %s found %s: %s", crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(), hex.EncodeToString(checkpoint), hex.EncodeToString(external2Signature), types.ErrInvalid),
		},
		{
			name: "bridger address not match",
			msg: &types.MsgOracleSetConfirm{
				Nonce:           nonce1OracleSet.Nonce,
				BridgerAddress:  suite.bridgers[1].String(),
				ExternalAddress: crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
				Signature:       hex.EncodeToString(external1Signature),
				ChainName:       suite.chainName,
			},
			err:       types.ErrInvalid,
			errReason: fmt.Sprintf("got %s, expected %s: %s", suite.bridgers[1].String(), suite.bridgers[0].String(), types.ErrInvalid),
		},
	}

	for _, testData := range errMsgData {
		_, err = suite.Handler()(suite.ctx, testData.msg)
		require.ErrorIs(suite.T(), err, testData.err, testData.name)
		require.EqualValues(suite.T(), err.Error(), testData.errReason, testData.name)
	}

	normalOracleSetConfirmMsg := &types.MsgOracleSetConfirm{
		Nonce:           nonce1OracleSet.Nonce,
		BridgerAddress:  suite.bridgers[0].String(),
		ExternalAddress: crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		Signature:       hex.EncodeToString(external1Signature),
		ChainName:       suite.chainName,
	}
	_, err = suite.Handler()(suite.ctx, normalOracleSetConfirmMsg)
	require.NoError(suite.T(), err)

	endBlockBeforeLatestOracleSet := suite.Keeper().GetLatestOracleSet(suite.ctx)
	require.NotNil(suite.T(), endBlockBeforeLatestOracleSet)
}

func (suite *MsgHandlerTestSuite) TestClaimWithOracleJailed() {
	normalMsg := &types.MsgCreateOracleBridger{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: suite.delegateAmount},
		ChainName:        suite.chainName,
	}
	_, err := suite.Handler()(suite.ctx, normalMsg)
	require.NoError(suite.T(), err)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	latestOracleSetNonce := suite.Keeper().GetLatestOracleSetNonce(suite.ctx)
	require.EqualValues(suite.T(), 1, latestOracleSetNonce)

	nonce1OracleSet := suite.Keeper().GetOracleSet(suite.ctx, latestOracleSetNonce)
	require.EqualValues(suite.T(), uint64(1), nonce1OracleSet.Nonce)
	require.EqualValues(suite.T(), uint64(2), nonce1OracleSet.Height)

	var gravityId string
	require.NotPanics(suite.T(), func() {
		gravityId = suite.Keeper().GetGravityID(suite.ctx)
	})
	require.EqualValues(suite.T(), "fx-bridge-bsc", gravityId)
	checkpoint, err := nonce1OracleSet.GetCheckpoint(gravityId)

	// oracle jailed!!!
	oracle, found := suite.Keeper().GetOracle(suite.ctx, suite.oracles[0])
	require.True(suite.T(), found)
	oracle.Jailed = true
	suite.Keeper().SetOracle(suite.ctx, oracle)

	external1Signature, err := types.NewEthereumSignature(checkpoint, suite.externals[0])
	require.NoError(suite.T(), err)

	normalOracleSetConfirmMsg := &types.MsgOracleSetConfirm{
		Nonce:           latestOracleSetNonce,
		BridgerAddress:  suite.bridgers[0].String(),
		ExternalAddress: crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		Signature:       hex.EncodeToString(external1Signature),
		ChainName:       suite.chainName,
	}
	_, err = suite.Handler()(suite.ctx, normalOracleSetConfirmMsg)
	require.Nil(suite.T(), err)
}

func (suite *MsgHandlerTestSuite) TestClaimTest() {
	normalMsg := &types.MsgCreateOracleBridger{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: suite.delegateAmount},
		ChainName:        suite.chainName,
	}
	_, err := suite.Handler()(suite.ctx, normalMsg)
	require.NoError(suite.T(), err)

	oracleLastEventNonce := suite.app.BscKeeper.GetLastEventNonceByOracle(suite.ctx, suite.oracles[0])
	require.EqualValues(suite.T(), 0, oracleLastEventNonce)

	errMsgDatas := []struct {
		name      string
		msg       *types.MsgBridgeTokenClaim
		err       error
		errReason string
	}{
		{
			name: "error oracleSet nonce: 2",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     2,
				BlockHeight:    1,
				TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
				Name:           "Pundix Token Purse",
				Symbol:         "PURSE",
				Decimals:       18,
				BridgerAddress: suite.bridgers[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      suite.chainName,
			},
			err:       types.ErrNonContiguousEventNonce,
			errReason: fmt.Sprintf("create attestation: got %v, expected %v: %s", 2, 1, types.ErrNonContiguousEventNonce),
		},
		{
			name: "error oracleSet nonce: 3",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     3,
				BlockHeight:    1,
				TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
				Name:           "Pundix Token Purse",
				Symbol:         "PURSE",
				Decimals:       18,
				BridgerAddress: suite.bridgers[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      suite.chainName,
			},
			err:       types.ErrNonContiguousEventNonce,
			errReason: fmt.Sprintf("create attestation: got %v, expected %v: %s", 3, 1, types.ErrNonContiguousEventNonce),
		},
		{
			name: "Normal claim msg: 1",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     1,
				BlockHeight:    1,
				TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
				Name:           "Pundix Token Purse",
				Symbol:         "PURSE",
				Decimals:       18,
				BridgerAddress: suite.bridgers[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      suite.chainName,
			},
			err:       nil,
			errReason: "",
		},
		{
			name: "error oracleSet nonce: 1",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     1,
				BlockHeight:    2,
				TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
				Name:           "Pundix Token Purse",
				Symbol:         "PURSE",
				Decimals:       18,
				BridgerAddress: suite.bridgers[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      suite.chainName,
			},
			err:       types.ErrNonContiguousEventNonce,
			errReason: fmt.Sprintf("create attestation: got %v, expected %v: %s", 1, 2, types.ErrNonContiguousEventNonce),
		},
		{
			name: "error oracleSet nonce: 3",
			msg: &types.MsgBridgeTokenClaim{
				EventNonce:     3,
				BlockHeight:    2,
				TokenContract:  "0x3f6795b8ABE0775a88973469909adE1405f7ac09",
				Name:           "Pundix Token Purse",
				Symbol:         "PURSE",
				Decimals:       18,
				BridgerAddress: suite.bridgers[0].String(),
				ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
				ChainName:      suite.chainName,
			},
			err:       types.ErrNonContiguousEventNonce,
			errReason: fmt.Sprintf("create attestation: got %v, expected %v: %s", 3, 2, types.ErrNonContiguousEventNonce),
		},
	}

	for _, testData := range errMsgDatas {
		_, err = suite.Handler()(suite.ctx, testData.msg)
		require.ErrorIs(suite.T(), err, testData.err, testData.name)
		if err == nil {
			continue
		}
		require.EqualValues(suite.T(), testData.errReason, err.Error(), testData.name)
	}

}

// Test RequestBatch baseFee
func (suite *MsgHandlerTestSuite) TestRequestBatchBaseFee() {

	endBlock := func() {
		crosschain.EndBlocker(suite.ctx, suite.Keeper())
	}

	// 1. First sets up a valid validator
	for i, oracle := range suite.oracles {
		normalMsg := &types.MsgCreateOracleBridger{
			OracleAddress:    oracle.String(),
			BridgerAddress:   suite.bridgers[i].String(),
			ExternalAddress:  crypto.PubkeyToAddress(suite.externals[i].PublicKey).Hex(),
			ValidatorAddress: suite.validator[i].String(),
			DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: suite.delegateAmount},
			ChainName:        suite.chainName,
		}
		_, err := suite.Handler()(suite.ctx, normalMsg)
		require.NoError(suite.T(), err)
	}

	endBlock()

	var externalOracleMembers []types.BridgeValidator
	for _, key := range suite.externals {
		externalOracleMembers = append(externalOracleMembers, types.BridgeValidator{
			Power:           100,
			ExternalAddress: crypto.PubkeyToAddress(key.PublicKey).Hex(),
		})
	}

	// 2. oracle update claim
	for i := range suite.oracles {
		normalMsg := &types.MsgOracleSetUpdatedClaim{
			EventNonce:     1,
			BlockHeight:    1,
			OracleSetNonce: 1,
			Members:        externalOracleMembers,
			BridgerAddress: suite.bridgers[i].String(),
			ChainName:      suite.chainName,
		}
		_, err := suite.Handler()(suite.ctx, normalMsg)
		require.NoError(suite.T(), err)
	}

	endBlock()

	// 3. add bridge token.
	sendToFxSendAddr := crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex()
	sendToFxReceiveAddr := suite.bridgers[0]
	sendToFxAmount := sdk.NewIntWithDecimal(1000, 18)
	sendToFxToken := "0x0000000000000000000000000000000000001000"

	for i, oracle := range suite.oracles {
		normalMsg := &types.MsgBridgeTokenClaim{
			EventNonce:     suite.Keeper().GetLastEventNonceByOracle(suite.ctx, oracle) + 1,
			BlockHeight:    1,
			TokenContract:  sendToFxToken,
			Name:           "BSC USDT",
			Symbol:         "USDT",
			Decimals:       18,
			BridgerAddress: suite.bridgers[i].String(),
			ChannelIbc:     "",
			ChainName:      suite.chainName,
		}
		_, err := suite.Handler()(suite.ctx, normalMsg)
		require.NoError(suite.T(), err)
	}

	endBlock()

	bridgeDenomData := suite.Keeper().GetBridgeTokenDenom(suite.ctx, sendToFxToken)
	require.NotNil(suite.T(), bridgeDenomData)
	tokenDenom := fmt.Sprintf("%s%s", suite.chainName, sendToFxToken)
	require.EqualValues(suite.T(), tokenDenom, bridgeDenomData.Denom)
	bridgeTokenData := suite.Keeper().GetDenomByBridgeToken(suite.ctx, tokenDenom)
	require.NotNil(suite.T(), bridgeTokenData)
	require.EqualValues(suite.T(), sendToFxToken, bridgeTokenData.Token)

	// 4. sendToFx.
	for i, oracle := range suite.oracles {
		normalMsg := &types.MsgSendToFxClaim{
			EventNonce:     suite.Keeper().GetLastEventNonceByOracle(suite.ctx, oracle) + 1,
			BlockHeight:    1,
			TokenContract:  sendToFxToken,
			Amount:         sendToFxAmount,
			Sender:         sendToFxSendAddr,
			Receiver:       sendToFxReceiveAddr.String(),
			TargetIbc:      "",
			BridgerAddress: suite.bridgers[i].String(),
			ChainName:      suite.chainName,
		}
		_, err := suite.Handler()(suite.ctx, normalMsg)
		require.NoError(suite.T(), err)
	}

	endBlock()

	balance := suite.app.BankKeeper.GetBalance(suite.ctx, sendToFxReceiveAddr, tokenDenom)
	require.NotNil(suite.T(), balance)
	require.EqualValues(suite.T(), balance.Denom, tokenDenom)
	require.True(suite.T(), balance.Amount.Equal(sendToFxAmount))

	sendToExternal := func(bridgeFees []sdk.Int) {
		for _, bridgeFee := range bridgeFees {
			sendToExternal := &types.MsgSendToExternal{
				Sender:    sendToFxReceiveAddr.String(),
				Dest:      sendToFxSendAddr,
				Amount:    sdk.NewCoin(tokenDenom, sdk.NewInt(3)),
				BridgeFee: sdk.NewCoin(tokenDenom, bridgeFee),
				ChainName: suite.chainName,
			}
			_, err := suite.Handler()(suite.ctx, sendToExternal)
			require.NoError(suite.T(), err)
		}
	}

	sendToExternal([]sdk.Int{sdk.NewInt(1), sdk.NewInt(2), sdk.NewInt(3)})
	usdtBatchFee := suite.Keeper().GetBatchFeesByTokenType(suite.ctx, sendToFxToken, 100, sdk.NewInt(0))
	require.EqualValues(suite.T(), sendToFxToken, usdtBatchFee.TokenContract)
	require.EqualValues(suite.T(), 3, usdtBatchFee.TotalTxs)
	require.EqualValues(suite.T(), sdk.NewInt(6), usdtBatchFee.TotalFees)

	getPointerInt := func(i sdk.Int) *sdk.Int {
		return &i
	}

	testCases := []struct {
		testName       string
		baseFee        *sdk.Int
		pass           bool
		expectTotalTxs uint64
		err            error
	}{
		{
			testName:       "Support - baseFee 1000",
			baseFee:        getPointerInt(sdk.NewInt(1000)),
			pass:           false,
			expectTotalTxs: 3,
			err:            types.ErrEmpty,
		},
		{
			testName:       "Support - baseFee 2",
			baseFee:        getPointerInt(sdk.NewInt(2)),
			pass:           true,
			expectTotalTxs: 1,
			err:            nil,
		},
		{
			testName:       "Support - baseFee 0",
			baseFee:        getPointerInt(sdk.NewInt(0)),
			pass:           true,
			expectTotalTxs: 0,
			err:            nil,
		},
	}

	for _, testCase := range testCases {
		suite.T().Run(testCase.testName, func(t *testing.T) {
			cacheCtx, _ := suite.ctx.CacheContext()
			_, err := suite.Handler()(cacheCtx, &types.MsgRequestBatch{
				Sender:     suite.bridgers[0].String(),
				Denom:      tokenDenom,
				MinimumFee: sdk.NewInt(1),
				FeeReceive: "0x0000000000000000000000000000000000000000",
				ChainName:  suite.chainName,
				BaseFee:    testCase.baseFee,
			})
			if testCase.pass {
				require.NoError(suite.T(), err)
				usdtBatchFee = suite.Keeper().GetBatchFeesByTokenType(cacheCtx, sendToFxToken, 100, sdk.NewInt(0))
				require.EqualValues(suite.T(), testCase.expectTotalTxs, usdtBatchFee.TotalTxs)
				return
			}

			require.NotNil(suite.T(), err)
			require.True(suite.T(), errors.As(err, &testCase.err))
			require.Equal(suite.T(), err, testCase.err)
		})
	}
}

func genEthKey(count int) []*ecdsa.PrivateKey {
	var ethKeys []*ecdsa.PrivateKey
	for i := 0; i < count; i++ {
		key, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		ethKeys = append(ethKeys, key)
	}
	return ethKeys
}
