package app

import (
	"fmt"
	"io"
	stdlog "log"
	"net/http"
	"os"
	"path/filepath"

	store "github.com/cosmos/cosmos-sdk/store/types"

	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"

	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/feegrant"

	othertypes "github.com/functionx/fx-core/x/other/types"

	"github.com/cosmos/cosmos-sdk/server/config"

	fxconfig "github.com/functionx/fx-core/server/config"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"

	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"

	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"

	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"

	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"

	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	ibc "github.com/cosmos/ibc-go/v3/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v3/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v3/modules/core/03-connection/types"
	porttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"

	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	// this line is used by starport scaffolding # stargate/app/moduleImport

	"github.com/functionx/fx-core/x/crosschain"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
	"github.com/functionx/fx-core/x/ibc/applications/transfer"
	ibctransferkeeper "github.com/functionx/fx-core/x/ibc/applications/transfer/keeper"
	ibctransfertypes "github.com/functionx/fx-core/x/ibc/applications/transfer/types"

	"github.com/functionx/fx-core/x/gravity"
	gravitykeeper "github.com/functionx/fx-core/x/gravity/keeper"
	gravitytypes "github.com/functionx/fx-core/x/gravity/types"

	"github.com/functionx/fx-core/x/bsc"
	bsctypes "github.com/functionx/fx-core/x/bsc/types"

	"github.com/functionx/fx-core/x/polygon"
	polygontypes "github.com/functionx/fx-core/x/polygon/types"

	crosschainkeeper "github.com/functionx/fx-core/x/crosschain/keeper"
	"github.com/functionx/fx-core/x/other"

	"github.com/functionx/fx-core/x/tron"
	tronkeeper "github.com/functionx/fx-core/x/tron/keeper"
	trontypes "github.com/functionx/fx-core/x/tron/types"

	"github.com/functionx/fx-core/x/evm"
	evmrest "github.com/functionx/fx-core/x/evm/client/rest"
	evmtypes "github.com/functionx/fx-core/x/evm/types"
	"github.com/functionx/fx-core/x/feemarket"
	feemarketkeeper "github.com/functionx/fx-core/x/feemarket/keeper"

	erc20client "github.com/functionx/fx-core/x/erc20/client"
	erc20keeper "github.com/functionx/fx-core/x/erc20/keeper"
	erc20types "github.com/functionx/fx-core/x/erc20/types"

	migratekeeper "github.com/functionx/fx-core/x/migrate/keeper"
	migratetypes "github.com/functionx/fx-core/x/migrate/types"

	"github.com/functionx/fx-core/app/ante"
	"github.com/functionx/fx-core/server"
	fxtypes "github.com/functionx/fx-core/types"
	"github.com/functionx/fx-core/x/erc20"
	evmkeeper "github.com/functionx/fx-core/x/evm/keeper"
	feemarkettypes "github.com/functionx/fx-core/x/feemarket/types"
	"github.com/functionx/fx-core/x/migrate"

	_ "github.com/functionx/fx-core/docs/statik"

	// Force-load the tracer engines to trigger registration due to Go-Ethereum v1.10.15 changes
	_ "github.com/ethereum/go-ethereum/eth/tracers/js"
	_ "github.com/ethereum/go-ethereum/eth/tracers/native"
)

func getGovProposalHandlers() []govclient.ProposalHandler {
	return []govclient.ProposalHandler{
		paramsclient.ProposalHandler,
		distrclient.ProposalHandler,
		upgradeclient.ProposalHandler,
		upgradeclient.CancelProposalHandler,
		erc20client.RegisterCoinProposalHandler,
		erc20client.ToggleTokenRelayProposalHandler,
		erc20client.RegisterERC20ProposalHandler,
	}
}

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(getGovProposalHandlers()...),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		// this line is used by starport scaffolding # stargate/app/moduleBasic
		other.AppModuleBasic{},
		gravity.AppModuleBasic{},
		crosschain.AppModuleBasic{},
		bsc.AppModuleBasic{},
		polygon.AppModuleBasic{},
		tron.AppModuleBasic{},
		evm.AppModuleBasic{},
		feemarket.AppModuleBasic{},
		erc20.AppModuleBasic{},
		migrate.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		// used for secure addition and subtraction of balance using module account
		gravitytypes.ModuleName: {authtypes.Minter, authtypes.Burner},
		bsctypes.ModuleName:     {authtypes.Minter, authtypes.Burner},
		polygontypes.ModuleName: {authtypes.Minter, authtypes.Burner},
		trontypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		evmtypes.ModuleName:     {authtypes.Minter, authtypes.Burner},
		erc20types.ModuleName:   {authtypes.Minter, authtypes.Burner},
	}
)

var (
	_ CosmosApp               = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	*baseapp.BaseApp

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper

	// IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCKeeper      *ibckeeper.Keeper
	EvidenceKeeper evidencekeeper.Keeper
	TransferKeeper ibctransferkeeper.Keeper
	FeeGrantKeeper feegrantkeeper.Keeper
	AuthzKeeper    authzkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper

	// gravity keepers
	GravityKeeper    gravitykeeper.Keeper
	CrosschainKeeper crosschainkeeper.RouterKeeper
	BscKeeper        crosschainkeeper.Keeper
	PolygonKeeper    crosschainkeeper.Keeper
	TronKeeper       crosschainkeeper.Keeper

	// ethermint keepers
	EvmKeeper       *evmkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper
	Erc20Keeper     erc20keeper.Keeper
	MigrateKeeper   migratekeeper.Keeper

	configurator module.Configurator

	// the module manager
	mm *module.Manager

	// simulation manager
	sm *module.SimulationManager
}

func init() {
	fxHome := os.ExpandEnv("$FX_HOME")
	if len(fxHome) > 0 {
		DefaultNodeHome = fxHome
		return
	}
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		stdlog.Println("Failed to get home dir %2", err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+fxtypes.Name)
}

func New(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	homePath string, invCheckPeriod uint, encodingConfig EncodingConfig,
	appOpts servertypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp),
) *App {

	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(fxtypes.Name, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		feegrant.StoreKey, authzkeeper.StoreKey,
		// this line is used by starport scaffolding # stargate/myApp/storeKey
		gravitytypes.StoreKey, bsctypes.StoreKey, polygontypes.StoreKey, trontypes.StoreKey,
		// ethermint keys
		evmtypes.StoreKey, feemarkettypes.StoreKey,
		// erc20 keys
		erc20types.StoreKey, migratetypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	myApp := &App{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	myApp.ParamsKeeper = initParamsKeeper(
		appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey],
	)

	// set the BaseApp's parameter store
	bApp.SetParamStore(
		myApp.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()),
	)

	// add capability keeper and ScopeToModule for ibc module
	myApp.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	// grant capabilities for the ibc and ibc-transfer modules
	scopedIBCKeeper := myApp.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := myApp.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	myApp.CapabilityKeeper.Seal()

	// add keepers
	myApp.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], myApp.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
	)
	myApp.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		myApp.BaseApp.MsgServiceRouter(),
	)
	myApp.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		keys[feegrant.StoreKey],
		myApp.AccountKeeper,
	)
	myApp.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], myApp.AccountKeeper, myApp.GetSubspace(banktypes.ModuleName), myApp.ModuleAccountAddrs(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], myApp.AccountKeeper, myApp.BankKeeper, myApp.GetSubspace(stakingtypes.ModuleName),
	)
	myApp.MintKeeper = mintkeeper.NewKeeper(
		appCodec, keys[minttypes.StoreKey], myApp.GetSubspace(minttypes.ModuleName), &stakingKeeper,
		myApp.AccountKeeper, myApp.BankKeeper, authtypes.FeeCollectorName,
	)
	myApp.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], myApp.GetSubspace(distrtypes.ModuleName), myApp.AccountKeeper, myApp.BankKeeper,
		&stakingKeeper, authtypes.FeeCollectorName, myApp.ModuleAccountAddrs(),
	)
	myApp.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, keys[slashingtypes.StoreKey], &stakingKeeper, myApp.GetSubspace(slashingtypes.ModuleName),
	)
	myApp.CrisisKeeper = crisiskeeper.NewKeeper(
		myApp.GetSubspace(crisistypes.ModuleName), invCheckPeriod, myApp.BankKeeper, authtypes.FeeCollectorName,
	)
	myApp.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, myApp.BaseApp)

	// Create IBC Keeper
	myApp.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey], myApp.GetSubspace(ibchost.ModuleName), stakingKeeper, myApp.UpgradeKeeper, scopedIBCKeeper,
	)

	// Create Transfer Keepers
	myApp.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], myApp.GetSubspace(ibctransfertypes.ModuleName),
		myApp.IBCKeeper.ChannelKeeper, &myApp.IBCKeeper.PortKeeper,
		myApp.AccountKeeper, myApp.BankKeeper, scopedTransferKeeper,
	)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], stakingKeeper, myApp.SlashingKeeper,
	)
	// If evidence needs to be handled for the myApp, set routes in router here and seal
	myApp.EvidenceKeeper = *evidenceKeeper

	// this line is used by starport scaffolding # stargate/myApp/keeperDefinition

	tracer := cast.ToString(appOpts.Get(server.EVMTracer))

	myApp.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec, keys[feemarkettypes.StoreKey], myApp.GetSubspace(feemarkettypes.ModuleName),
	)

	myApp.EvmKeeper = evmkeeper.NewKeeper(
		appCodec, keys[evmtypes.StoreKey], tkeys[evmtypes.TransientKey], myApp.GetSubspace(evmtypes.ModuleName),
		myApp.AccountKeeper, myApp.BankKeeper, stakingKeeper, myApp.FeeMarketKeeper, tracer)

	myApp.Erc20Keeper = erc20keeper.NewKeeper(
		keys[erc20types.StoreKey], appCodec, myApp.GetSubspace(erc20types.ModuleName),
		myApp.AccountKeeper, myApp.BankKeeper, myApp.EvmKeeper,
		&myApp.TransferKeeper, myApp.IBCKeeper.ChannelKeeper)

	myApp.GravityKeeper = gravitykeeper.NewKeeper(
		appCodec, keys[gravitytypes.StoreKey], myApp.GetSubspace(gravitytypes.ModuleName),
		stakingKeeper, myApp.BankKeeper, myApp.AccountKeeper, myApp.SlashingKeeper,
		myApp.TransferKeeper, myApp.IBCKeeper.ChannelKeeper, myApp.Erc20Keeper,
	)

	// init cross chain module
	myApp.BscKeeper = crosschainkeeper.NewKeeper(
		appCodec, bsctypes.ModuleName, keys[bsctypes.StoreKey], myApp.GetSubspace(bsctypes.ModuleName),
		myApp.BankKeeper, myApp.AccountKeeper, myApp.TransferKeeper, myApp.IBCKeeper.ChannelKeeper, myApp.Erc20Keeper)

	myApp.PolygonKeeper = crosschainkeeper.NewKeeper(
		appCodec, polygontypes.ModuleName, keys[polygontypes.StoreKey], myApp.GetSubspace(polygontypes.ModuleName),
		myApp.BankKeeper, myApp.AccountKeeper, myApp.TransferKeeper, myApp.IBCKeeper.ChannelKeeper, myApp.Erc20Keeper)

	myApp.TronKeeper = crosschainkeeper.NewKeeper(
		appCodec, trontypes.ModuleName, keys[trontypes.StoreKey], myApp.GetSubspace(trontypes.ModuleName),
		myApp.BankKeeper, myApp.AccountKeeper, myApp.TransferKeeper, myApp.IBCKeeper.ChannelKeeper, myApp.Erc20Keeper)

	// add cross-chain router
	crosschainRouter := crosschainkeeper.NewRouter()
	crosschainRouter.
		AddRoute(bsctypes.ModuleName, &crosschainkeeper.ModuleHandler{
			QueryServer: myApp.BscKeeper,
			MsgServer:   crosschainkeeper.NewMsgServerImpl(myApp.BscKeeper),
		}).
		AddRoute(polygontypes.ModuleName, &crosschainkeeper.ModuleHandler{
			QueryServer: myApp.PolygonKeeper,
			MsgServer:   crosschainkeeper.NewMsgServerImpl(myApp.PolygonKeeper),
		}).
		AddRoute(trontypes.ModuleName, &crosschainkeeper.ModuleHandler{
			QueryServer: myApp.TronKeeper,
			MsgServer:   tronkeeper.NewMsgServerImpl(myApp.TronKeeper),
		})

	myApp.CrosschainKeeper = crosschainkeeper.NewRouterKeeper(crosschainRouter)

	// register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(myApp.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(myApp.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(myApp.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(myApp.IBCKeeper.ClientKeeper)).
		AddRoute(crosschaintypes.RouterKey, crosschain.NewCrossChainProposalHandler(myApp.CrosschainKeeper)).
		AddRoute(erc20types.RouterKey, erc20.NewErc20ProposalHandler(&myApp.Erc20Keeper))

	myApp.GovKeeper = govkeeper.NewKeeper(
		appCodec, keys[govtypes.StoreKey], myApp.GetSubspace(govtypes.ModuleName), myApp.AccountKeeper, myApp.BankKeeper, &stakingKeeper, govRouter,
	)

	erc20TransferRouter := erc20types.NewRouter()
	erc20TransferRouter.AddRoute(gravitytypes.ModuleName, myApp.GravityKeeper)
	erc20TransferRouter.AddRoute(bsctypes.ModuleName, myApp.BscKeeper)
	erc20TransferRouter.AddRoute(polygontypes.ModuleName, myApp.PolygonKeeper)
	erc20TransferRouter.AddRoute(trontypes.ModuleName, myApp.TronKeeper)
	myApp.Erc20Keeper.SetRouter(erc20TransferRouter)

	ibcTransferRouter := ibctransfertypes.NewRouter()
	ibcTransferRouter.AddRoute(gravitytypes.ModuleName, myApp.GravityKeeper)
	ibcTransferRouter.AddRoute(bsctypes.ModuleName, myApp.BscKeeper)
	ibcTransferRouter.AddRoute(polygontypes.ModuleName, myApp.PolygonKeeper)
	ibcTransferRouter.AddRoute(trontypes.ModuleName, myApp.TronKeeper)
	ibcTransferRouter.AddRoute(erc20types.ModuleName, myApp.Erc20Keeper)
	myApp.TransferKeeper.SetRouter(ibcTransferRouter)
	myApp.TransferKeeper.SetRefundHook(myApp.Erc20Keeper)
	transferModule := transfer.NewAppModule(myApp.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(myApp.TransferKeeper)
	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferIBCModule)
	// this line is used by starport scaffolding # ibc/myApp/router
	myApp.IBCKeeper.SetRouter(ibcRouter)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	myApp.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			myApp.DistrKeeper.Hooks(),
			myApp.SlashingKeeper.Hooks(),
			myApp.GravityKeeper.Hooks(),
		),
	)
	myApp.EvmKeeper.SetHooks(evmkeeper.NewMultiEvmHooks(myApp.Erc20Keeper.Hooks()))

	myApp.MigrateKeeper = migratekeeper.NewKeeper(appCodec, keys[migratetypes.StoreKey], myApp.AccountKeeper)
	bankMigrate := migratekeeper.NewBankMigrate(myApp.BankKeeper)
	distrStakingMigrate := migratekeeper.NewDistrStakingMigrate(keys[distrtypes.StoreKey], keys[stakingtypes.StoreKey], myApp.StakingKeeper)
	govMigrate := migratekeeper.NewGovMigrate(keys[govtypes.StoreKey], myApp.GovKeeper)
	myApp.MigrateKeeper.SetMigrateI(bankMigrate, distrStakingMigrate, govMigrate)
	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	var skipGenesisInvariants = cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	myApp.mm = module.NewManager(
		genutil.NewAppModule(
			myApp.AccountKeeper, myApp.StakingKeeper, myApp.BaseApp.DeliverTx, encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, myApp.AccountKeeper, nil),
		vesting.NewAppModule(myApp.AccountKeeper, myApp.BankKeeper),
		bank.NewAppModule(appCodec, myApp.BankKeeper, myApp.AccountKeeper),
		capability.NewAppModule(appCodec, *myApp.CapabilityKeeper),
		crisis.NewAppModule(&myApp.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, myApp.GovKeeper, myApp.AccountKeeper, myApp.BankKeeper),
		mint.NewAppModule(appCodec, myApp.MintKeeper, myApp.AccountKeeper),
		slashing.NewAppModule(appCodec, myApp.SlashingKeeper, myApp.AccountKeeper, myApp.BankKeeper, myApp.StakingKeeper),
		distr.NewAppModule(appCodec, myApp.DistrKeeper, myApp.AccountKeeper, myApp.BankKeeper, myApp.StakingKeeper),
		staking.NewAppModule(appCodec, myApp.StakingKeeper, myApp.AccountKeeper, myApp.BankKeeper),
		upgrade.NewAppModule(myApp.UpgradeKeeper),
		evidence.NewAppModule(myApp.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, myApp.AccountKeeper, myApp.BankKeeper, myApp.FeeGrantKeeper, myApp.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, myApp.AuthzKeeper, myApp.AccountKeeper, myApp.BankKeeper, myApp.interfaceRegistry),
		ibc.NewAppModule(myApp.IBCKeeper),
		params.NewAppModule(myApp.ParamsKeeper),
		// this line is used by starport scaffolding # stargate/myApp/appModule
		gravity.NewAppModule(myApp.GravityKeeper, myApp.BankKeeper),
		other.NewAppModule(appCodec),
		crosschain.NewAppModuleByRouter(myApp.CrosschainKeeper),
		bsc.NewAppModule(myApp.BscKeeper, myApp.BankKeeper),
		polygon.NewAppModule(myApp.PolygonKeeper, myApp.BankKeeper),
		tron.NewAppModule(myApp.TronKeeper, myApp.BankKeeper),
		// Ethermint app modules
		evm.NewAppModule(myApp.EvmKeeper, myApp.AccountKeeper),
		feemarket.NewAppModule(myApp.FeeMarketKeeper),
		erc20.NewAppModule(myApp.Erc20Keeper, myApp.AccountKeeper),
		migrate.NewAppModule(myApp.MigrateKeeper),
		transferModule,
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	// NOTE: capability module's beginblocker must come before any modules using capabilities (e.g. IBC)
	myApp.mm.SetOrderBeginBlockers(
		// upgrades should be run first
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		//liquiditytypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		//icatypes.ModuleName,
		//routertypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		capabilitytypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,

		othertypes.ModuleName,
		gravitytypes.ModuleName,
		crosschaintypes.ModuleName,
		bsctypes.ModuleName,
		trontypes.ModuleName,
		polygontypes.ModuleName,
		migratetypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		erc20types.ModuleName,
	)

	myApp.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		//liquiditytypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		//icatypes.ModuleName,
		//routertypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,

		othertypes.ModuleName,
		gravitytypes.ModuleName,
		crosschaintypes.ModuleName,
		bsctypes.ModuleName,
		trontypes.ModuleName,
		polygontypes.ModuleName,
		migratetypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		erc20types.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: The genutils module must also occur after auth so that it can access the params from auth.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	myApp.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		//icatypes.ModuleName,
		evidencetypes.ModuleName,
		//liquiditytypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		authtypes.ModuleName,
		genutiltypes.ModuleName,
		//routertypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,

		othertypes.ModuleName,
		gravitytypes.ModuleName,
		crosschaintypes.ModuleName,
		bsctypes.ModuleName,
		trontypes.ModuleName,
		polygontypes.ModuleName,
		migratetypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		erc20types.ModuleName,
	)

	myApp.mm.RegisterInvariants(&myApp.CrisisKeeper)
	myApp.mm.RegisterRoutes(myApp.Router(), myApp.QueryRouter(), encodingConfig.Amino)

	myApp.configurator = module.NewConfigurator(myApp.appCodec, myApp.MsgServiceRouter(), myApp.GRPCQueryRouter())
	myApp.mm.RegisterServices(myApp.configurator)

	// initialize stores
	myApp.MountKVStores(keys)
	myApp.MountTransientStores(tkeys)
	myApp.MountMemoryStores(memKeys)

	maxGasWanted := cast.ToUint64(appOpts.Get(server.EVMMaxTxGasWanted))
	anteOptions := ante.HandlerOptions{
		AccountKeeper:        myApp.AccountKeeper,
		BankKeeper:           myApp.BankKeeper,
		EvmKeeper:            myApp.EvmKeeper,
		IbcKeeper:            myApp.IBCKeeper,
		SignModeHandler:      encodingConfig.TxConfig.SignModeHandler(),
		SigGasConsumer:       ante.DefaultSigVerificationGasConsumer,
		MaxTxGasWanted:       maxGasWanted,
		BypassMinFeeMsgTypes: cast.ToStringSlice(appOpts.Get(fxconfig.BypassMinFeeMsgTypesKey)),
	}

	if err := anteOptions.Validate(); err != nil {
		panic(fmt.Errorf("failed to ante options validate: %s", err))
	}

	myApp.SetAnteHandler(ante.NewAnteHandler(anteOptions))
	myApp.SetInitChainer(myApp.InitChainer)
	myApp.SetBeginBlocker(myApp.BeginBlocker)
	myApp.SetEndBlocker(myApp.EndBlocker)

	myApp.UpgradeKeeper.SetUpgradeHandler("v2", func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// set max expected block time parameter. Replace the default with your expected value
		// https://github.com/cosmos/ibc-go/blob/release/v1.0.x/docs/ibc/proto-docs.md#params-2
		myApp.IBCKeeper.ConnectionKeeper.SetParams(ctx, ibcconnectiontypes.DefaultParams())

		fromVM = map[string]uint64{
			// other modules
			ibchost.ModuleName:        1,
			evmtypes.ModuleName:       1,
			feemarkettypes.ModuleName: 1,
			erc20types.ModuleName:     1,
			migratetypes.ModuleName:   1,
		}

		ctx.Logger().Info("start to run module migrations...")

		return myApp.mm.RunMigrations(ctx, myApp.configurator, fromVM)
	})

	// TODO need fix
	//rootmulti.AddIgnoreCommitKey(fxtypes.CrossChainSupportBscBlock(), bsctypes.StoreKey)
	//rootmulti.AddIgnoreCommitKey(fxtypes.CrossChainSupportPolygonAndTronBlock(), polygontypes.StoreKey, trontypes.StoreKey)
	//rootmulti.AddIgnoreCommitKey(fxtypes.EvmV0SupportBlock(), evmtypesv0.StoreKey, feemarkettypesv0.StoreKey)
	//rootmulti.AddIgnoreCommitKey(fxtypes.EvmV1SupportBlock(), evmtypes.StoreKey, feemarkettypes.StoreKey, erc20types.StoreKey, migratetypes.StoreKey)
	//govtypes.SetEGFProposalSupportBlock(fxtypes.EvmV1SupportBlock())

	upgradeInfo, err := myApp.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == "v2" && !myApp.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := store.StoreUpgrades{
			Added: []string{
				evmtypes.StoreKey,
				feemarkettypes.StoreKey,
				erc20types.StoreKey,
				migratetypes.StoreKey,
			},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		myApp.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}

	if loadLatest {
		if err := myApp.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	myApp.ScopedIBCKeeper = scopedIBCKeeper
	myApp.ScopedTransferKeeper = scopedTransferKeeper

	return myApp
}

// Name returns the name of the App
func (app *App) Name() string {
	return app.BaseApp.Name()
}

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())

	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *App) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// LegacyAmino NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns InterfaceRegistry
func (app *App) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *App) GetMemKey(storeKey string) *sdk.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// SimulationManager implements the SimulationApp interface
func (app *App) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	// Register evm tx routes.
	evmrest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	// Register legacy tx routes.
	authrest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(apiSvr.Router)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	// this line is used by starport scaffolding # stargate/app/paramSubspace
	paramsKeeper.Subspace(gravitytypes.ModuleName)
	paramsKeeper.Subspace(bsctypes.ModuleName)
	paramsKeeper.Subspace(polygontypes.ModuleName)
	paramsKeeper.Subspace(trontypes.ModuleName)

	// ethermint subspaces
	paramsKeeper.Subspace(evmtypes.ModuleName)
	paramsKeeper.Subspace(feemarkettypes.ModuleName)
	paramsKeeper.Subspace(erc20types.ModuleName)
	return paramsKeeper
}
