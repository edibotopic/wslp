package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"wslp/internal/wsl"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered WSL distros",
	Long:  `Lists all WSL distributions registered on the Windows host.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Finding registered distros...")

		distros, err := wsl.ListDistros(context.Background())

		if err != nil {
			fmt.Println(err)
			return
		}

		// print the number of distros
		fmt.Println(len(distros), "distros are registered:")

		// print the list of distros
		for i := range distros {
			fmt.Println(distros[i].Name)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
