package cmd

import (
	"testing"
)

func TestServeCommand(t *testing.T) {
	t.Run("command metadata", func(t *testing.T) {
		if RootCmd == nil {
			t.Fatal("RootCmd is nil")
		}

		serveCmd, _, err := RootCmd.Find([]string{"serve"})
		if err != nil {
			t.Fatalf("serve command not found: %v", err)
		}

		if serveCmd.Use != "serve" {
			t.Errorf("expected Use='serve', got '%s'", serveCmd.Use)
		}

		if serveCmd.Short == "" {
			t.Error("Short description is empty")
		}

		if serveCmd.Long == "" {
			t.Error("Long description is empty")
		}
	})

	t.Run("port flag exists and defaults to 8080", func(t *testing.T) {
		serveCmd, _, err := RootCmd.Find([]string{"serve"})
		if err != nil {
			t.Fatalf("serve command not found: %v", err)
		}

		portFlag := serveCmd.Flags().Lookup("port")
		if portFlag == nil {
			t.Fatal("port flag not found")
		}

		// Check default value
		if portFlag.DefValue != "8080" {
			t.Errorf("expected port default to be '8080', got '%s'", portFlag.DefValue)
		}
	})

	t.Run("port flag has short form", func(t *testing.T) {
		serveCmd, _, err := RootCmd.Find([]string{"serve"})
		if err != nil {
			t.Fatalf("serve command not found: %v", err)
		}

		portFlag := serveCmd.Flags().Lookup("port")
		if portFlag == nil {
			t.Fatal("port flag not found")
		}

		if portFlag.Shorthand != "p" {
			t.Errorf("expected port short flag to be 'p', got '%s'", portFlag.Shorthand)
		}
	})
}
