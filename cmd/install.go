package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	gowsl "github.com/ubuntu/gowsl"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install WSL distros",
	Long:  `Install one or more WSL distros`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: No distros specified")
			return
		}

		// TODO: try to install distros concurrently
		for _, distro := range args {
			fmt.Printf("Installing %s...\n", distro)
			if err := gowsl.Install(context.Background(), distro); err != nil {
				fmt.Printf("Error installing %s: %v\n", distro, err)
				continue
			}
			fmt.Printf("Successfully installed %s\n", distro)

			d := gowsl.NewDistro(context.Background(), distro)
			registered, err := d.IsRegistered()
			if err != nil {
				fmt.Printf("Error checking registration: %v\n", err)
			} else if registered {
				fmt.Printf("%s was downloaded in the modern format and is already registered\n", distro)
				fmt.Printf("Launch with wsl -d %s", distro)
			} else {
				fmt.Printf("%s was downloaded in the classic format and must be registered\n", distro)
				fmt.Printf("Register with wsl --register %s", distro)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
