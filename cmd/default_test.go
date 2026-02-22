package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// TestDefaultCommand tests the default command
func TestDefaultCommand(t *testing.T) {
	// Test 1: Check the command exists and has correct metadata
	// This test works on any OS since it only checks command structure
	t.Run("command metadata", func(t *testing.T) {
		if RootCmd == nil {
			t.Fatal("RootCmd is nil")
		}

		defaultCmd, _, err := RootCmd.Find([]string{"default"})
		if err != nil {
			t.Fatalf("default command not found: %v", err)
		}

		if defaultCmd.Use != "default" {
			t.Errorf("expected Use='default', got '%s'", defaultCmd.Use)
		}

		if defaultCmd.Short == "" {
			t.Error("Short description is empty")
		}

		if defaultCmd.Long == "" {
			t.Error("Long description is empty")
		}
	})

	// Test 2: Execute the command and verify output
	// This test requires Windows and is skipped on Linux/Mac
	t.Run("command output", func(t *testing.T) {
		skipIfNotWindows(t)

		buf := new(bytes.Buffer)
		
		// Create a fresh command instance for this test
		cmd := RootCmd
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs([]string{"default", "show"})
		defer cmd.SetArgs([]string{}) // Reset after test

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}

		output := buf.String()

		// Check for expected output patterns
		expectedPhrases := []string{
			"The default WSL distro is:",
		}

		for _, phrase := range expectedPhrases {
			if !strings.Contains(output, phrase) {
				t.Errorf("expected output to contain '%s', got:\n%s", phrase, output)
			}
		}
	})
}
