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

		// TODO: install distros concurrently
		// TODO: investigate why classic distros (old format) don't install
		// TODO: investigate if custom name for distro can be passed
		for _, distro := range args {
			fmt.Printf("Installing %s...\n", distro)
			if err := gowsl.Install(context.Background(), distro); err != nil {
				fmt.Printf("Error installing %s: %v\n", distro, err)
				continue
			}
			fmt.Printf("Successfully installed %s\n", distro)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
