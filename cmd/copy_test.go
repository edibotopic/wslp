package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

type mockCopier struct {
	shouldFail bool
	registered map[string]bool // track which names are registered
}

func (m *mockCopier) IsRegistered(ctx context.Context, name string) (bool, error) {
	if m.shouldFail {
		return false, errors.New("mock error")
	}
	if m.registered != nil {
		return m.registered[name], nil
	}
	// Default: source exists, new name doesn't
	return name == "Ubuntu", nil
}

func (m *mockCopier) Export(ctx context.Context, distroName, outputPath string) error {
	if m.shouldFail {
		return errors.New("export failed")
	}
	return nil
}

func (m *mockCopier) Import(ctx context.Context, newName, tarPath, installDir string) error {
	if m.shouldFail {
		return errors.New("import failed")
	}
	return nil
}

func TestCopyCommand(t *testing.T) {
	t.Run("command metadata", func(t *testing.T) {
		if RootCmd == nil {
			t.Fatal("RootCmd is nil")
		}

		copyCmd, _, err := RootCmd.Find([]string{"copy"})
		if err != nil {
			t.Fatalf("copy command not found: %v", err)
		}

		if copyCmd.Use != "copy <source> <new-name>" {
			t.Errorf("expected Use='copy <source> <new-name>', got '%s'", copyCmd.Use)
		}

		if copyCmd.Short == "" {
			t.Error("Short description is empty")
		}

		if copyCmd.Long == "" {
			t.Error("Long description is empty")
		}
	})

	t.Run("copy distro success", func(t *testing.T) {
		mock := &mockCopier{}
		out := new(bytes.Buffer)

		err := CopyDistroCmd(context.Background(), mock, out, "Ubuntu", "UbuntuCopy", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "Copying Ubuntu to UbuntuCopy") {
			t.Errorf("expected copy message in output, got:\n%s", output)
		}

		if !strings.Contains(output, "✓") {
			t.Errorf("expected success indicator in output, got:\n%s", output)
		}
	})

	t.Run("returns error when copy fails", func(t *testing.T) {
		mock := &mockCopier{shouldFail: true}
		out := new(bytes.Buffer)

		err := CopyDistroCmd(context.Background(), mock, out, "Ubuntu", "UbuntuCopy", "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "copy failed") {
			t.Errorf("expected 'copy failed' error, got: %v", err)
		}
	})

	t.Run("prints failure message on copy failure", func(t *testing.T) {
		mock := &mockCopier{shouldFail: true}
		out := new(bytes.Buffer)

		CopyDistroCmd(context.Background(), mock, out, "Ubuntu", "UbuntuCopy", "")

		output := out.String()
		if !strings.Contains(output, "✗") {
			t.Errorf("expected failure indicator in output, got:\n%s", output)
		}
	})

	t.Run("flag defaults", func(t *testing.T) {
		copyCmd, _, err := RootCmd.Find([]string{"copy"})
		if err != nil {
			t.Fatalf("copy command not found: %v", err)
		}

		installDirFlag := copyCmd.Flags().Lookup("install-dir")
		if installDirFlag == nil {
			t.Fatal("install-dir flag not found")
		}
	})
}
