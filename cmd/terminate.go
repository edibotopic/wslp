package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"wslp/internal/wsl"
)

func init() {
	RootCmd.AddCommand(newTerminateCmd())
}

func newTerminateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "terminate <distro> [distro...]",
		Aliases: []string{"stop", "kill"},
		Short:   "Terminate one or more running WSL distributions",
		Long: `Terminate (stop) one or more running WSL distributions.

This is useful before performing operations like backups, or to free up system resources.
Terminating a distro will stop all processes running in that distribution.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTerminate(cmd.OutOrStdout(), args)
		},
	}

	return cmd
}

func runTerminate(w io.Writer, distros []string) error {
	ctx := context.Background()

	terminator := wsl.RealTerminator{}
	results := wsl.TerminateDistros(ctx, terminator, distros)

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
		fmt.Fprintf(w, "\nSuccessfully terminated %d/%d distribution(s)\n", successCount, len(results))
	}

	if successCount < len(results) {
		return fmt.Errorf("some terminations failed")
	}

	return nil
}
