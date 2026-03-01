package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"wslp/internal/wsl"
)

func init() {
	RootCmd.AddCommand(newRenameCmd())
}

func newRenameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename <old-name> <new-name>",
		Short: "Rename a WSL distribution",
		Long: `Rename a WSL distribution by modifying the Windows Registry.

This command updates the distribution name in the registry. After renaming,
you should restart WSL for the changes to take effect:

    wsl --shutdown

The rename operation:
- Validates the old distro exists
- Checks the new name doesn't conflict with existing distros
- Updates the registry entry directly (fast, no export/import needed)`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRename(cmd.OutOrStdout(), args[0], args[1])
		},
	}

	return cmd
}

func runRename(w io.Writer, oldName, newName string) error {
	ctx := context.Background()

	renamer := wsl.RealRenamer{}
	result := wsl.RenameDistro(ctx, renamer, oldName, newName)

	if result.Success {
		fmt.Fprintf(w, "✓ %s\n", result.Message)
		fmt.Fprintf(w, "\nTo apply changes, run:\n")
		fmt.Fprintf(w, "  wsl --shutdown\n")
		return nil
	} else {
		fmt.Fprintf(w, "✗ %s\n", result.Message)
		return fmt.Errorf("rename failed")
	}
}
