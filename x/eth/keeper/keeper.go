package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	crosschainkeeper "github.com/functionx/fx-core/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/x/crosschain/types"
)

type Keeper struct {
	crosschainkeeper.Keeper
}

func NewKeeper(cdc codec.BinaryCodec, moduleName string, storeKey sdk.StoreKey, paramSpace paramtypes.Subspace,
	bankKeeper crosschaintypes.BankKeeper, ibcTransferKeeper crosschaintypes.IBCTransferKeeper,
	channelKeeper crosschaintypes.IBCChannelKeeper, erc20Keeper crosschaintypes.Erc20Keeper) Keeper {
	return Keeper{
		Keeper: crosschainkeeper.NewKeeper(
			cdc, moduleName, storeKey, paramSpace, bankKeeper, ibcTransferKeeper, channelKeeper, erc20Keeper,
		),
	}
}