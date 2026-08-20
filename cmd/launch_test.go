package cmd

import (
	"testing"
)

func TestLaunchCommand(t *testing.T) {
	t.Run("command metadata", func(t *testing.T) {
		if RootCmd == nil {
			t.Fatal("RootCmd is nil")
		}

		launchCmd, _, err := RootCmd.Find([]string{"launch"})
		if err != nil {
			t.Fatalf("launch command not found: %v", err)
		}

		if launchCmd.Use != "launch <distro>" {
			t.Errorf("expected Use='launch <distro>', got '%s'", launchCmd.Use)
		}

		if launchCmd.Short == "" {
			t.Error("Short description is empty")
		}

		if launchCmd.Long == "" {
			t.Error("Long description is empty")
		}
	})

	t.Run("launch has aliases", func(t *testing.T) {
		launchCmd, _, err := RootCmd.Find([]string{"launch"})
		if err != nil {
			t.Fatalf("launch command not found: %v", err)
		}

		hasStart := false
		hasOpen := false
		for _, alias := range launchCmd.Aliases {
			if alias == "start" {
				hasStart = true
			}
			if alias == "open" {
				hasOpen = true
			}
		}

		if !hasStart {
			t.Error("expected 'start' alias")
		}
		if !hasOpen {
			t.Error("expected 'open' alias")
		}
	})

	t.Run("launch command exists and is runnable", func(t *testing.T) {
		launchCmd, _, err := RootCmd.Find([]string{"launch"})
		if err != nil {
			t.Fatalf("launch command not found: %v", err)
		}

		// Verify it's a cobra command with Run/RunE set
		if launchCmd.RunE == nil && launchCmd.Run == nil {
			t.Error("launch command has no Run or RunE function")
		}
	})
}
