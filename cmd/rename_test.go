package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

type mockRenamer struct {
	shouldFail bool
	registered map[string]bool // track which names are registered
}

func (m *mockRenamer) IsRegistered(ctx context.Context, name string) (bool, error) {
	if m.shouldFail {
		return false, errors.New("mock error")
	}
	if m.registered != nil {
		return m.registered[name], nil
	}
	// Default: old name exists, new name doesn't
	return name == "Ubuntu", nil
}

func (m *mockRenamer) GetDistroGUID(ctx context.Context, name string) (string, error) {
	if m.shouldFail {
		return "", errors.New("failed to get GUID")
	}
	return "12345678-1234-1234-1234-123456789012", nil
}

func (m *mockRenamer) RenameInRegistry(guid, newName string) error {
	if m.shouldFail {
		return errors.New("failed to update registry")
	}
	return nil
}

func TestRenameCommand(t *testing.T) {
	t.Run("command metadata", func(t *testing.T) {
		if RootCmd == nil {
			t.Fatal("RootCmd is nil")
		}

		renameCmd, _, err := RootCmd.Find([]string{"rename"})
		if err != nil {
			t.Fatalf("rename command not found: %v", err)
		}

		if renameCmd.Use != "rename <old-name> <new-name>" {
			t.Errorf("expected Use='rename <old-name> <new-name>', got '%s'", renameCmd.Use)
		}

		if renameCmd.Short == "" {
			t.Error("Short description is empty")
		}

		if renameCmd.Long == "" {
			t.Error("Long description is empty")
		}
	})

	t.Run("rename distro success", func(t *testing.T) {
		mock := &mockRenamer{}
		out := new(bytes.Buffer)

		err := RenameDistroCmd(context.Background(), mock, out, "Ubuntu", "MyUbuntu")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "✓") {
			t.Errorf("expected success indicator in output, got:\n%s", output)
		}

		if !strings.Contains(output, "wsl --shutdown") {
			t.Errorf("expected shutdown instruction in output, got:\n%s", output)
		}
	})

	t.Run("returns error when rename fails", func(t *testing.T) {
		mock := &mockRenamer{shouldFail: true}
		out := new(bytes.Buffer)

		err := RenameDistroCmd(context.Background(), mock, out, "Ubuntu", "MyUbuntu")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "rename failed") {
			t.Errorf("expected 'rename failed' error, got: %v", err)
		}
	})

	t.Run("prints failure message on rename failure", func(t *testing.T) {
		mock := &mockRenamer{shouldFail: true}
		out := new(bytes.Buffer)

		RenameDistroCmd(context.Background(), mock, out, "Ubuntu", "MyUbuntu")

		output := out.String()
		if !strings.Contains(output, "✗") {
			t.Errorf("expected failure indicator in output, got:\n%s", output)
		}

		if strings.Contains(output, "wsl --shutdown") {
			t.Errorf("should not include shutdown instruction on failure, got:\n%s", output)
		}
	})
}
