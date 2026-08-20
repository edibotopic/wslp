package cmd

import (
	"testing"
)

func TestRootCommand(t *testing.T) {
	t.Run("RootCmd exists and has proper metadata", func(t *testing.T) {
		if RootCmd == nil {
			t.Fatal("RootCmd is nil")
		}

		if RootCmd.Use != "wslp" {
			t.Errorf("expected Use='wslp', got '%s'", RootCmd.Use)
		}

		if RootCmd.Short == "" {
			t.Error("Short description is empty")
		}

		if RootCmd.Long == "" {
			t.Error("Long description is empty")
		}
	})

	t.Run("common subcommands are registered", func(t *testing.T) {
		expectedCommands := []string{"list", "default", "backup", "copy", "terminate", "rename", "unregister", "install", "launch", "serve"}

		for _, cmd := range expectedCommands {
			found, _, err := RootCmd.Find([]string{cmd})
			if err != nil {
				t.Errorf("command '%s' not found: %v", cmd, err)
			}

			if found == nil {
				t.Errorf("command '%s' is nil", cmd)
			}
		}
	})
}
