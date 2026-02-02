package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	gowsl "github.com/ubuntu/gowsl"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered WSL distros",
	Long:  `Lists all WSL distributions registered on the Windows host.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Finding registered distros...")

		distroList, err := gowsl.RegisteredDistros(context.Background())

		if err != nil {
			fmt.Println(err)
		}

		// print the number of distros
		fmt.Println(len(distroList), "distros are registered:")

		// print the list of distros
		for i := range distroList {
			fmt.Println(distroList[i].Name())
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
