package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	gowsl "github.com/ubuntu/gowsl"
)

// Installer interface for mocking the install operation
// We don't want to actually install something during tests
type Installer interface {
	Install(ctx context.Context, distro string) error
}

// RealInstaller uses the actual gowsl library
type RealInstaller struct{}

func (r RealInstaller) Install(ctx context.Context, distro string) error {
	return gowsl.Install(ctx, distro)
}

// Default installer (real implementation)
var currentInstaller Installer = RealInstaller{}

// InstallDistros has the core logic for installing distros
func InstallDistros(ctx context.Context, installer Installer, out io.Writer, distros []string) error {
	if len(distros) == 0 {
		fmt.Fprintln(out, "Error: No distros specified")
		return nil
	}

	// TODO: try to install distros concurrently
	for _, distro := range distros {
		fmt.Fprintf(out, "Installing %s...\n", distro)
		if err := installer.Install(ctx, distro); err != nil {
			fmt.Fprintf(out, "Error installing %s: %v\n", distro, err)
			continue
		}
		fmt.Fprintf(out, "Successfully installed %s\n", distro)

		d := gowsl.NewDistro(ctx, distro)
		registered, err := d.IsRegistered()
		if err != nil {
			fmt.Fprintf(out, "Error checking registration: %v\n", err)
		} else if registered {
			fmt.Fprintf(out, "%s was downloaded in the modern format and is already registered\n", distro)
			fmt.Fprintf(out, "Launch with wsl -d %s\n", distro)
		} else {
			fmt.Fprintf(out, "%s was downloaded in the classic format and must be registered\n", distro)
			fmt.Fprintf(out, "Register with wsl --register %s\n", distro)
		}
	}

	return nil
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install WSL distros",
	Long:  `Install one or more WSL distros`,
	Run: func(cmd *cobra.Command, args []string) {
		InstallDistros(context.Background(), currentInstaller, cmd.OutOrStdout(), args)
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
