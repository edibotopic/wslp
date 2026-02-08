package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// TestInstallCommand tests the install command
func TestInstallCommand(t *testing.T) {
	// Test 1: Check the command exists and has correct metadata
	// This test works on any OS since it only checks command structure
	t.Run("command metadata", func(t *testing.T) {
		if RootCmd == nil {
			t.Fatal("RootCmd is nil")
		}

		installCmd, _, err := RootCmd.Find([]string{"install"})
		if err != nil {
			t.Fatalf("install command not found: %v", err)
		}

		if installCmd.Use != "install" {
			t.Errorf("expected Use='install', got '%s'", installCmd.Use)
		}

		if installCmd.Short == "" {
			t.Error("Short description is empty")
		}

		if installCmd.Long == "" {
			t.Error("Long description is empty")
		}
	})

	// Test 2: Execute the command and verify output
	// This test requires Windows and is skipped on Linux/Mac
	t.Run("command output", func(t *testing.T) {
		skipIfNotWindows(t)

		buf := new(bytes.Buffer)
		RootCmd.SetOut(buf)
		RootCmd.SetErr(buf)
		RootCmd.SetArgs([]string{"install", "Ubuntu"})

		err := RootCmd.Execute()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}

		output := buf.String()

		// Check for expected output patterns
		expectedPhrases := []string{
			"Installing",
			"Ubuntu",
		}

		for _, phrase := range expectedPhrases {
			if !strings.Contains(output, phrase) {
				t.Errorf("expected output to contain '%s', got:\n%s", phrase, output)
			}
		}
	})
}
