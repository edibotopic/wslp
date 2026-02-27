package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"wslp/internal/wsl"
)

func ListDistros(ctx context.Context, l wsl.Lister, out io.Writer) error {
	fmt.Fprintln(out, "Finding registered distros...")

	distros, err := wsl.ListDistros(ctx, l)
	if err != nil {
		fmt.Fprintln(out, err)
		return err
	}

	fmt.Fprintf(out, "%d distros are registered:\n", len(distros))
	for _, d := range distros {
		fmt.Fprintln(out, d.Name)
	}

	return nil
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered WSL distros",
	Long:  `Lists all WSL distributions registered on the Windows host.`,
	Run: func(cmd *cobra.Command, args []string) {
		ListDistros(context.Background(), wsl.RealLister{}, cmd.OutOrStdout())
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
