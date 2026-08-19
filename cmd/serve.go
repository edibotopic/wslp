package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"

	"wslp/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP API server",
	Long:  `Starts an HTTP server that exposes WSL operations via REST API for the Flutter GUI.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")

		s := server.NewServer(port)

		errCh := make(chan error, 1)
		go func() {
			errCh <- s.Start()
		}()

		// Also allow the server to be stopped via the /api/shutdown
		// endpoint (used by the GUI), or via Ctrl+C on the console.
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		select {
		case err := <-errCh:
			if err != nil {
				fmt.Printf("Server error: %v\n", err)
			}
		case <-sigCh:
			fmt.Println("Shutting down server...")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := s.Shutdown(ctx); err != nil {
				fmt.Printf("Error during shutdown: %v\n", err)
			}
			<-errCh
		}

		fmt.Println("Server stopped")
	},
}

func init() {
	serveCmd.Flags().StringP("port", "p", "8080", "Port to run the server on")
	RootCmd.AddCommand(serveCmd)
}
