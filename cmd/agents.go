package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed agents_snippet.md
var agentsSnippet string

func newAgentsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "agents",
		Short: "Print recommended AGENTS.md instructions",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprint(cmd.OutOrStdout(), agentsSnippet)
			return err
		},
	}
}
