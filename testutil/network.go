package testutil

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	hd2 "github.com/evmos/ethermint/crypto/hd"
	dbm "github.com/tendermint/tm-db"

	"github.com/functionx/fx-core/v3/app"
	"github.com/functionx/fx-core/v3/app/helpers"
	"github.com/functionx/fx-core/v3/testutil/network"
	fxtypes "github.com/functionx/fx-core/v3/types"
)

func DefNoSupplyGenesisState(cdc codec.Codec) app.GenesisState {
	genesisState := app.NewDefAppGenesisByDenom(fxtypes.DefaultDenom, cdc)
	bankState := banktypes.DefaultGenesisState()
	bankState.DenomMetadata = []banktypes.Metadata{fxtypes.GetFXMetaData(fxtypes.DefaultDenom)}
	genesisState[banktypes.ModuleName] = cdc.MustMarshalJSON(bankState)
	return genesisState
}

// DefaultNetworkConfig returns a sane default configuration suitable for nearly all
// testing requirements.
func DefaultNetworkConfig() network.Config {
	encCfg := app.MakeEncodingConfig()

	return network.Config{
		Codec:             encCfg.Marshaler,
		TxConfig:          encCfg.TxConfig,
		LegacyAmino:       encCfg.Amino,
		InterfaceRegistry: encCfg.InterfaceRegistry,
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor: func(val network.Validator) servertypes.Application {
			return app.New(
				val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
				encCfg,
				helpers.EmptyAppOptions{},
				baseapp.SetPruning(storetypes.NewPruningOptionsFromString(val.AppConfig.Pruning)),
				baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
			)
		},
		GenesisState:    DefNoSupplyGenesisState(encCfg.Marshaler),
		TimeoutCommit:   500 * time.Millisecond,
		ChainID:         fxtypes.MainnetChainId,
		NumValidators:   4,
		BondDenom:       fxtypes.DefaultDenom,
		MinGasPrices:    fmt.Sprintf("4000000000000%s", fxtypes.DefaultDenom),
		AccountTokens:   sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction),
		StakingTokens:   sdk.TokensFromConsensusPower(500, sdk.DefaultPowerReduction),
		BondedTokens:    sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction),
		PruningStrategy: storetypes.PruningOptionNothing,
		CleanupDir:      true,
		SigningAlgo:     string(hd.Secp256k1Type),
		KeyringOptions: []keyring.Option{
			hd2.EthSecp256k1Option(),
		},
		PrintMnemonic: false,
	}
}
