package wsl

import (
	"bytes"
	"testing"
)

func TestPrintInstallResults(t *testing.T) {
	t.Run("prints success message for registered distro (modern format)", func(t *testing.T) {
		var buf bytes.Buffer
		results := []InstallResult{
			{
				Distro:     "Ubuntu",
				Success:    true,
				Registered: true,
				Message:    "Successfully installed (modern format, already registered)",
			},
		}

		PrintInstallResults(&buf, results)
		output := buf.String()

		if !bytes.Contains(buf.Bytes(), []byte("Successfully installed Ubuntu")) {
			t.Errorf("expected success message, got: %s", output)
		}
		if !bytes.Contains(buf.Bytes(), []byte("modern format")) {
			t.Errorf("expected modern format message, got: %s", output)
		}
		if !bytes.Contains(buf.Bytes(), []byte("already registered")) {
			t.Errorf("expected registered message, got: %s", output)
		}
		if !bytes.Contains(buf.Bytes(), []byte("wsl -d Ubuntu")) {
			t.Errorf("expected launch command, got: %s", output)
		}
	})

	t.Run("prints success and registration hint for unregistered distro (classic format)", func(t *testing.T) {
		var buf bytes.Buffer
		results := []InstallResult{
			{
				Distro:     "Debian",
				Success:    true,
				Registered: false,
				Message:    "Successfully installed (classic format, must be registered with: wsl --register Debian)",
			},
		}

		PrintInstallResults(&buf, results)
		output := buf.String()

		if !bytes.Contains(buf.Bytes(), []byte("Successfully installed Debian")) {
			t.Errorf("expected success message, got: %s", output)
		}
		if !bytes.Contains(buf.Bytes(), []byte("classic format")) {
			t.Errorf("expected classic format message, got: %s", output)
		}
		if !bytes.Contains(buf.Bytes(), []byte("wsl --register")) {
			t.Errorf("expected register hint, got: %s", output)
		}
		if !bytes.Contains(buf.Bytes(), []byte("Debian")) {
			t.Errorf("expected distro name in hint, got: %s", output)
		}
	})

	t.Run("prints error message for failed installation", func(t *testing.T) {
		var buf bytes.Buffer
		results := []InstallResult{
			{
				Distro:  "AlpineLinux",
				Success: false,
				Message: "Download failed: connection timeout",
			},
		}

		PrintInstallResults(&buf, results)
		output := buf.String()

		if !bytes.Contains(buf.Bytes(), []byte("Error installing AlpineLinux")) {
			t.Errorf("expected error message, got: %s", output)
		}
		if !bytes.Contains(buf.Bytes(), []byte("connection timeout")) {
			t.Errorf("expected error details, got: %s", output)
		}
	})

	t.Run("prints multiple results correctly", func(t *testing.T) {
		var buf bytes.Buffer
		results := []InstallResult{
			{
				Distro:     "Ubuntu",
				Success:    true,
				Registered: true,
				Message:    "Successfully installed (modern format, already registered)",
			},
			{
				Distro:  "InvalidDistro",
				Success: false,
				Message: "Distribution not found",
			},
			{
				Distro:     "Fedora",
				Success:    true,
				Registered: false,
				Message:    "Successfully installed (classic format, must be registered)",
			},
		}

		PrintInstallResults(&buf, results)
		output := buf.String()

		if !bytes.Contains(buf.Bytes(), []byte("Ubuntu")) {
			t.Errorf("expected Ubuntu in output, got: %s", output)
		}
		if !bytes.Contains(buf.Bytes(), []byte("InvalidDistro")) {
			t.Errorf("expected InvalidDistro in output, got: %s", output)
		}
		if !bytes.Contains(buf.Bytes(), []byte("Fedora")) {
			t.Errorf("expected Fedora in output, got: %s", output)
		}
	})

	t.Run("handles empty results slice", func(t *testing.T) {
		var buf bytes.Buffer
		results := []InstallResult{}

		PrintInstallResults(&buf, results)
		output := buf.String()

		if len(output) != 0 {
			t.Errorf("expected no output for empty results, got: %s", output)
		}
	})

	t.Run("prints error for successful install but registration check error", func(t *testing.T) {
		var buf bytes.Buffer
		results := []InstallResult{
			{
				Distro:  "CentOS",
				Success: true,
				Message: "Successfully installed (Error checking registration: registry error)",
			},
		}

		PrintInstallResults(&buf, results)
		output := buf.String()

		if !bytes.Contains(buf.Bytes(), []byte("CentOS")) {
			t.Errorf("expected distro name in output, got: %s", output)
		}
		if !bytes.Contains(buf.Bytes(), []byte("Successfully installed")) {
			t.Errorf("expected success message, got: %s", output)
		}
	})
}
