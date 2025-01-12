package v6

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v6/app/keepers"
	crosschainkeeper "github.com/functionx/fx-core/v6/x/crosschain/keeper"
	fxgovtypes "github.com/functionx/fx-core/v6/x/gov/types"
	layer2types "github.com/functionx/fx-core/v6/x/layer2/types"
	migratekeeper "github.com/functionx/fx-core/v6/x/migrate/keeper"
	fxstakingkeeper "github.com/functionx/fx-core/v6/x/staking/keeper"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	app *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		cacheCtx, commit := ctx.CacheContext()

		if err := UpdateParams(cacheCtx, app); err != nil {
			return nil, err
		}

		MigrateMetadata(cacheCtx, app.BankKeeper)
		MigrateLayer2Module(cacheCtx, app.Layer2Keeper)
		ExportCommunityPool(cacheCtx, app.DistrKeeper, app.BankKeeper)

		ctx.Logger().Info("start to run v6 migrations...", "module", "upgrade")
		toVM, err := mm.RunMigrations(cacheCtx, configurator, fromVM)
		if err != nil {
			return fromVM, err
		}

		commit()
		ctx.Logger().Info("Upgrade complete")
		return toVM, nil
	}
}

func UpdateParams(cacheCtx sdk.Context, app *keepers.AppKeepers) error {
	mintParams := app.MintKeeper.GetParams(cacheCtx)
	mintParams.InflationMax = sdk.ZeroDec()
	mintParams.InflationMin = sdk.ZeroDec()
	if err := mintParams.Validate(); err != nil {
		return err
	}
	app.MintKeeper.SetParams(cacheCtx, mintParams)

	distrParams := app.DistrKeeper.GetParams(cacheCtx)
	distrParams.CommunityTax = sdk.ZeroDec()
	distrParams.BaseProposerReward = sdk.ZeroDec()
	distrParams.BonusProposerReward = sdk.ZeroDec()
	if err := distrParams.ValidateBasic(); err != nil {
		return err
	}
	app.DistrKeeper.SetParams(cacheCtx, distrParams)

	stakingParams := app.StakingKeeper.GetParams(cacheCtx)
	stakingParams.UnbondingTime = 1
	if err := stakingParams.Validate(); err != nil {
		return err
	}
	app.StakingKeeper.SetParams(cacheCtx, stakingParams)

	govTallyParams := app.GovKeeper.GetTallyParams(cacheCtx)
	govTallyParams.Quorum = sdk.OneDec().String()        // 100%
	govTallyParams.Threshold = sdk.OneDec().String()     // 100%
	govTallyParams.VetoThreshold = sdk.OneDec().String() // 100%
	app.GovKeeper.SetTallyParams(cacheCtx, govTallyParams)

	app.GovKeeper.IterateParams(cacheCtx, func(param *fxgovtypes.Params) (stop bool) {
		param.Quorum = sdk.OneDec().String()        // 100%
		param.Threshold = sdk.OneDec().String()     // 100%
		param.VetoThreshold = sdk.OneDec().String() // 100%
		if err := param.ValidateBasic(); err != nil {
			panic(err)
		}
		if err := app.GovKeeper.SetParams(cacheCtx, param); err != nil {
			panic(err)
		}
		return false
	})
	return nil
}

func ExportCommunityPool(ctx sdk.Context, distrKeeper distrkeeper.Keeper, bankKeeper bankkeeper.Keeper) sdk.Coins {
	feePool := distrKeeper.GetFeePool(ctx)
	truncatedCoins, changeCoins := feePool.CommunityPool.TruncateDecimal()
	feePool.CommunityPool = changeCoins
	distrKeeper.SetFeePool(ctx, feePool)

	if err := bankKeeper.SendCoinsFromModuleToModule(ctx, distrtypes.ModuleName, govtypes.ModuleName, truncatedCoins); err != nil {
		panic(err)
	}
	if err := bankKeeper.BurnCoins(ctx, govtypes.ModuleName, truncatedCoins); err != nil {
		panic(err)
	}
	ctx.Logger().Info("export community pool", "coins", truncatedCoins.String())
	return truncatedCoins
}

func MigrateMetadata(ctx sdk.Context, bankKeeper bankkeeper.Keeper) {
	bankKeeper.IterateAllDenomMetaData(ctx, func(metadata banktypes.Metadata) bool {
		address, ok := Layer2GenesisTokenAddress[metadata.Symbol]
		if !ok {
			return false
		}
		if len(metadata.DenomUnits) > 0 {
			metadata.DenomUnits[0].Aliases = append(metadata.DenomUnits[0].Aliases,
				fmt.Sprintf("%s%s", layer2types.ModuleName, address))
			bankKeeper.SetDenomMetaData(ctx, metadata)
		}
		return false
	})
}

func MigrateLayer2Module(ctx sdk.Context, layer2CrossChainKeeper crosschainkeeper.Keeper) {
	for _, address := range Layer2GenesisTokenAddress {
		fxTokenDenom := fmt.Sprintf("%s%s", layer2types.ModuleName, address)
		layer2CrossChainKeeper.AddBridgeToken(ctx, address, fxTokenDenom)
	}
}

func AutoUndelegate(ctx sdk.Context, stakingKeeper fxstakingkeeper.Keeper) []stakingtypes.Delegation {
	var delegations []stakingtypes.Delegation
	stakingKeeper.IterateAllDelegations(ctx, func(delegation stakingtypes.Delegation) (stop bool) {
		delegations = append(delegations, delegation)
		delegator := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)
		valAddress, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		if delegator.Equals(valAddress) {
			return false
		}
		if _, err := stakingKeeper.Undelegate(ctx, delegator, valAddress, delegation.Shares); err != nil {
			panic(err)
		}
		return false
	})
	return delegations
}

func ExportDelegate(ctx sdk.Context, migrateKeeper migratekeeper.Keeper, delegations []stakingtypes.Delegation) []stakingtypes.Delegation {
	for i := 0; i < len(delegations); i++ {
		delegation := delegations[i]
		delegator := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)
		if !migrateKeeper.HasMigratedDirectionTo(ctx, common.BytesToAddress(delegator.Bytes())) {
			delegations = append(delegations[:i], delegations[i+1:]...)
			i--
			continue
		}
	}
	return delegations
}
