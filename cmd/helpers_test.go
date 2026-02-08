package cmd

import (
	"runtime"
	"testing"
)

// skipIfNotWindows skips the test if not running on Windows
func skipIfNotWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skipf("Skipping test - requires Windows, running on %s", runtime.GOOS)
	}
}
