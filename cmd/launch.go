package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"wslp/internal/wsl"
)

func init() {
	RootCmd.AddCommand(newLaunchCmd())
}

func newLaunchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "launch <distro>",
		Aliases: []string{"start", "open"},
		Short:   "Launch an interactive shell for a WSL distribution",
		Long: `Launch an interactive shell for the specified WSL distribution.

This opens the default shell for the distro in the current terminal window.
The command will block until you exit the shell.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return wsl.LaunchInteractive(context.Background(), args[0])
		},
	}

	return cmd
}
