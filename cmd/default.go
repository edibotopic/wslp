package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"wslp/internal/wsl"
)

var defaultCmd = &cobra.Command{
	Use:   "default",
	Short: "Manage the default WSL distro",
	Long:  `Manage the default Linux distro, including showing it (done) and changing it (todo).`,
}

func ShowDefault(ctx context.Context, g wsl.DefaultGetter, out io.Writer) error {
	fmt.Fprintln(out, "The default WSL distro is:")

	name, err := wsl.GetDefaultDistro(ctx, g)
	if err != nil {
		fmt.Fprintln(out, err)
		return err
	}

	fmt.Fprintln(out, name)
	return nil
}

var defaultShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the default distro",
	Long:  `Prints the default WSL distribution on the Windows host.`,
	Run: func(cmd *cobra.Command, args []string) {
		ShowDefault(context.Background(), wsl.RealDefaultGetter{}, cmd.OutOrStdout())
	},
}

// TODO
var defaultChangeCmd = &cobra.Command{
	Use:   "change [distroName]",
	Short: "Change the default distro (STUB: not implemented)",
	Long:  `Changes the default WSL distro.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This subcommand is not implemented (sorry!).")
	},
}

func init() {
	// Add the top-level default command to root
	RootCmd.AddCommand(defaultCmd)

	// Add subcommands to the default command
	defaultCmd.AddCommand(defaultShowCmd)
	defaultCmd.AddCommand(defaultChangeCmd)
}
