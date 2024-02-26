package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/evmos/ethermint/x/bridge/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(_ string) *cobra.Command {
	// Group bridge queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdListSignerGroup())
	cmd.AddCommand(CmdShowSignerGroup())
	cmd.AddCommand(CmdListCallerGroup())
	cmd.AddCommand(CmdShowCallerGroup())
	cmd.AddCommand(CmdListDeposit())
	cmd.AddCommand(CmdShowDeposit())
	cmd.AddCommand(CmdListWithdraw())
	cmd.AddCommand(CmdListWithdrawByStatus())
	cmd.AddCommand(CmdShowWithdraw())
	cmd.AddCommand(CmdListRollupTx())
	cmd.AddCommand(CmdShowRollupTx())
	// this line is used by starport scaffolding # 1

	return cmd
}