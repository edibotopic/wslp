package wsl

import (
	"testing"
)

func TestIsWindowsTerminalAvailable(t *testing.T) {
	t.Run("returns a bool without panicking", func(t *testing.T) {
		// Test that the function returns a bool and doesn't panic.
		// The actual result (true/false) is environment-dependent and depends on
		// whether Windows Terminal is installed on the machine running the test,
		// so we just verify it returns a bool successfully.
		result := isWindowsTerminalAvailable()

		// Verify it's a bool (this will always pass, but documents the contract)
		if _, ok := interface{}(result).(bool); !ok {
			t.Errorf("expected bool, got %T", result)
		}
	})

	t.Run("function completes without panic", func(t *testing.T) {
		// This test simply ensures the function completes execution
		// and doesn't panic, which is the primary concern for exec.LookPath
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("isWindowsTerminalAvailable panicked: %v", r)
			}
		}()

		_ = isWindowsTerminalAvailable()
	})
}
