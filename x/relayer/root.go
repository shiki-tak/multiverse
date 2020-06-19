package relayer

import (
	"github.com/shiki-tak/connect/x/relayer/config"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "relayer",
		Aliases: []string{"cfg"},
		Short:   "manage configuration file",
	}

	cmd.AddCommand(
		config.ConfigCmd(),
	)

	return cmd
}
