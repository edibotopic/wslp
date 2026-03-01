package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"wslp/internal/wsl"
)

func init() {
	RootCmd.AddCommand(newUnregisterCmd())
}

func newUnregisterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unregister <distro> [distro...]",
		Aliases: []string{"delete", "remove"},
		Short:   "Unregister one or more WSL distributions",
		Long: `Unregister (delete) one or more WSL distributions.

WARNING: This will permanently delete the distribution and all its data.
Make sure to backup any important data before unregistering.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUnregister(cmd.OutOrStdout(), args)
		},
	}

	return cmd
}

func runUnregister(w io.Writer, distros []string) error {
	ctx := context.Background()

	unregisterer := wsl.RealUnregisterer{}
	results := wsl.UnregisterDistros(ctx, unregisterer, distros)

	// Print results
	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
			fmt.Fprintf(w, "✓ %s: %s\n", result.Distro, result.Message)
		} else {
			fmt.Fprintf(w, "✗ %s: %s\n", result.Distro, result.Message)
		}
	}

	if successCount > 0 {
		fmt.Fprintf(w, "\nSuccessfully unregistered %d/%d distribution(s)\n", successCount, len(results))
	}

	if successCount < len(results) {
		return fmt.Errorf("some unregistrations failed")
	}

	return nil
}
