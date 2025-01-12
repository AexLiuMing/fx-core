package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/functionx/fx-core/v6/x/crosschain/client/cli"
	"github.com/functionx/fx-core/v6/x/layer2/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the layer2 module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(cli.GetQuerySubCmds(types.ModuleName)...)
	return cmd
}
