package cmd

import (
	"fmt"

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
		if err := s.Start(); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	},
}

func init() {
	serveCmd.Flags().StringP("port", "p", "8080", "Port to run the server on")
	RootCmd.AddCommand(serveCmd)
}
