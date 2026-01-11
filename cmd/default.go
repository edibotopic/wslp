package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ubuntu/gowsl"
)

var defaultCmd = &cobra.Command{
	Use:   "default",
	Short: "Manage the default WSL distro",
	Long:  `Set the default Linux distro, change the default (TODO), and switch between WSL 1 and 2 (TODO).`,
}

var defaultShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the default distro",
	Long:  `Prints the defaults WSL distribution on the Windows host.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("The default WSL distro is:")

		defaultDistro, _, err := gowsl.DefaultDistro(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(defaultDistro.Name())
	},
}

var defaultChangeCmd = &cobra.Command{
	Use:   "change [distroName]",
	Short: "Change the default distro (TODO)",
	Long:  `A longer description that spans multiple lines and likely contains examples.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This subcommand is not implemented.")
	},
}

// defaultCmd represents the default command
var defaultChangeVersionCmd = &cobra.Command{
	Use:   "change [wslVersion]",
	Short: "Change the default WSL version (1 or 2)",
	Long:  `Sets the default WSL version on the Windows host..`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This subcommand is not implemented.")
	},
}

func init() {
	// Add the top-level default command to root
	rootCmd.AddCommand(defaultCmd)

	// Add subcommands to the default command
	defaultCmd.AddCommand(defaultShowCmd)
	defaultCmd.AddCommand(defaultChangeCmd)
	defaultCmd.AddCommand(defaultChangeVersionCmd)
}
