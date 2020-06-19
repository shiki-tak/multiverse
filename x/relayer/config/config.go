package config

import (
	relayercmd "github.com/shiki-tak/relayer/cmd"
	"github.com/spf13/cobra"
)

func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "manage configuration file",
	}

	cmd.AddCommand(
		relayercmd.ConfigShowCmd(),
		relayercmd.ConfigInitCmd(),
		relayercmd.ConfigAddDirCmd(),
	)

	return cmd
}

func InitConfig(cmd *cobra.Command) error {
	return relayercmd.InitConfig(cmd)
}
