package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// MockInstaller is a fake installer for testing only
type MockInstaller struct {
	installedDistros []string
	shouldFail       bool
}

func (m *MockInstaller) Install(ctx context.Context, distro string) error {
	m.installedDistros = append(m.installedDistros, distro)
	if m.shouldFail {
		return context.DeadlineExceeded // Simulate an error
	}
	return nil
}

// TestInstallCommand tests the install command
func TestInstallCommand(t *testing.T) {
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

	t.Run("no distros specified", func(t *testing.T) {
		mock := &MockInstaller{}
		out := new(bytes.Buffer)

		err := InstallDistros(context.Background(), mock, out, []string{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "No distros specified") {
			t.Errorf("expected error message about no distros, got: %s", output)
		}

		if len(mock.installedDistros) > 0 {
			t.Errorf("expected no distros installed, got: %v", mock.installedDistros)
		}
	})

	t.Run("install single distro", func(t *testing.T) {
		mock := &MockInstaller{}
		out := new(bytes.Buffer)

		err := InstallDistros(context.Background(), mock, out, []string{"Ubuntu"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()

		expectedPhrases := []string{
			"Installing Ubuntu",
		}

		for _, phrase := range expectedPhrases {
			if !strings.Contains(output, phrase) {
				t.Errorf("expected output to contain '%s', got:\n%s", phrase, output)
			}
		}

		if len(mock.installedDistros) != 1 || mock.installedDistros[0] != "Ubuntu" {
			t.Errorf("expected Ubuntu to be installed, got: %v", mock.installedDistros)
		}
	})

	t.Run("install multiple distros", func(t *testing.T) {
		mock := &MockInstaller{}
		out := new(bytes.Buffer)

		err := InstallDistros(context.Background(), mock, out, []string{"Ubuntu", "Debian", "Kali"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()

		for _, distro := range []string{"Ubuntu", "Debian", "Kali"} {
			if !strings.Contains(output, "Installing "+distro) {
				t.Errorf("expected output to mention installing %s, got:\n%s", distro, output)
			}
		}

		if len(mock.installedDistros) != 3 {
			t.Errorf("expected 3 distros installed, got: %v", mock.installedDistros)
		}
	})

	t.Run("handle install failure", func(t *testing.T) {
		mock := &MockInstaller{shouldFail: true}
		out := new(bytes.Buffer)

		err := InstallDistros(context.Background(), mock, out, []string{"Ubuntu"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()

		if !strings.Contains(output, "Error installing") {
			t.Errorf("expected error message in output, got:\n%s", output)
		}
	})
}
