package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// TestListCommand tests the list command
func TestListCommand(t *testing.T) {
	// Test 1: Check the command exists and has correct metadata
	// This test works on any OS since it only checks command structure
	t.Run("command metadata", func(t *testing.T) {
		if RootCmd == nil {
			t.Fatal("RootCmd is nil")
		}

		listCmd, _, err := RootCmd.Find([]string{"list"})
		if err != nil {
			t.Fatalf("list command not found: %v", err)
		}

		if listCmd.Use != "list" {
			t.Errorf("expected Use='list', got '%s'", listCmd.Use)
		}

		if listCmd.Short == "" {
			t.Error("Short description is empty")
		}

		if listCmd.Long == "" {
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
		cmd.SetArgs([]string{"list"})
		defer cmd.SetArgs([]string{}) // Reset after test

		err := cmd.Execute()
		if err != nil {
			t.Fatalf("command failed: %v", err)
		}

		output := buf.String()

		// Check for expected output patterns
		expectedPhrases := []string{
			"Finding registered distros",
			"distros are registered",
		}

		for _, phrase := range expectedPhrases {
			if !strings.Contains(output, phrase) {
				t.Errorf("expected output to contain '%s', got:\n%s", phrase, output)
			}
		}
	})
}
