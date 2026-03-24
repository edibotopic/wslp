package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"wslp/internal/wsl"
)

func init() {
	RootCmd.AddCommand(newCopyCmd())
}

func newCopyCmd() *cobra.Command {
	var installDir string

	cmd := &cobra.Command{
		Use:   "copy <source> <new-name>",
		Short: "Copy a WSL distribution under a new name",
		Long: `Copy a WSL distribution by exporting it and importing it under a new name.

The new distribution is stored in %USERPROFILE%\WSLCopies\<new-name> by default.
You can override this with the --install-dir flag.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCopy(cmd.OutOrStdout(), args[0], args[1], installDir)
		},
	}

	cmd.Flags().StringVarP(&installDir, "install-dir", "d", "", "Directory to store the new distro's virtual disk (overrides default)")

	return cmd
}

func runCopy(w io.Writer, source, newName, installDir string) error {
	ctx := context.Background()

	fmt.Fprintf(w, "Copying %s to %s...\n", source, newName)

	copier := wsl.RealCopier{}
	result := wsl.CopyDistro(ctx, copier, source, newName, installDir)

	if result.Success {
		fmt.Fprintf(w, "✓ %s\n", result.Message)
	} else {
		fmt.Fprintf(w, "✗ %s\n", result.Message)
		return fmt.Errorf("copy failed")
	}

	return nil
}
