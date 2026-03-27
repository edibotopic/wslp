package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"wslp/internal/config"
	"wslp/internal/wsl"
)

// InstallDistros has the core logic for installing distros
func InstallDistros(ctx context.Context, out io.Writer, distros []string) {
	results := wsl.InstallDistros(ctx, distros, false)
	wsl.PrintInstallResults(out, results)
}

// InstallDistrosConcurrent installs distros concurrently using a semaphore
func InstallDistrosConcurrent(ctx context.Context, out io.Writer, distros []string) {
	results := wsl.InstallDistros(ctx, distros, true)
	wsl.PrintInstallResults(out, results)
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install <distro> [distro...]",
	Short: "Install WSL distros",
	Long:  `Install one or more WSL distros`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "Error: No distros specified")
			return
		}
		concurrent, _ := cmd.Flags().GetBool("experimental-concurrent")
		if concurrent {
			fmt.Fprintf(cmd.OutOrStdout(), "experimental: installing distros concurrently (max %d at a time)\n", config.GetMaxConcurrentInstalls())
			InstallDistrosConcurrent(context.Background(), cmd.OutOrStdout(), args)
		} else {
			InstallDistros(context.Background(), cmd.OutOrStdout(), args)
		}
	},
}

func init() {
	installCmd.Flags().Bool("experimental-concurrent", false, "experimental: install distros concurrently")
	RootCmd.AddCommand(installCmd)
}
