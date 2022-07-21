package types

import (
	crosschaintypes "github.com/functionx/fx-core/v2/x/crosschain/types"
)

func DefaultGenesisState() *crosschaintypes.GenesisState {
	params := crosschaintypes.DefaultParams()
	params.GravityId = "fx-bsc-bridge"
	return &crosschaintypes.GenesisState{
		Params: params,
	}
}
