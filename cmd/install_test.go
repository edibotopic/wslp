package cmd

import (
	"bytes"
	"context"
	"testing"
)

func TestInstallCommand(t *testing.T) {
	t.Run("command metadata", func(t *testing.T) {
		if RootCmd == nil {
			t.Fatal("RootCmd is nil")
		}

		installCmd, _, err := RootCmd.Find([]string{"install"})
		if err != nil {
			t.Fatalf("install command not found: %v", err)
		}

		if installCmd.Use != "install <distro> [distro...]" {
			t.Errorf("expected Use='install <distro> [distro...]', got '%s'", installCmd.Use)
		}

		if installCmd.Short == "" {
			t.Error("Short description is empty")
		}

		if installCmd.Long == "" {
			t.Error("Long description is empty")
		}
	})

	t.Run("concurrent flag exists and defaults to false", func(t *testing.T) {
		installCmd, _, err := RootCmd.Find([]string{"install"})
		if err != nil {
			t.Fatalf("install command not found: %v", err)
		}

		concurrentFlag := installCmd.Flags().Lookup("experimental-concurrent")
		if concurrentFlag == nil {
			t.Fatal("experimental-concurrent flag not found")
		}

		// Check default value
		if concurrentFlag.DefValue != "false" {
			t.Errorf("expected experimental-concurrent default to be 'false', got '%s'", concurrentFlag.DefValue)
		}
	})

	t.Run("InstallDistros exports a function", func(t *testing.T) {
		// This verifies the exported function exists and can be called
		// We use an empty list here to avoid actual installation attempts
		out := new(bytes.Buffer)
		InstallDistros(context.Background(), out, []string{})
		// Just verify it doesn't panic
	})

	t.Run("InstallDistrosConcurrent exports a function", func(t *testing.T) {
		// This verifies the exported function exists and can be called
		// We use an empty list here to avoid actual installation attempts
		out := new(bytes.Buffer)
		InstallDistrosConcurrent(context.Background(), out, []string{})
		// Just verify it doesn't panic
	})
}
