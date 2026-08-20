package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

type mockBackuper struct {
	shouldFail bool
	failOn     string // specific distro to fail on
}

func (m *mockBackuper) IsRegistered(ctx context.Context, name string) (bool, error) {
	if m.shouldFail && m.failOn == name {
		return false, errors.New("mock error")
	}
	return true, nil
}

func (m *mockBackuper) Export(ctx context.Context, distroName, outputPath string) error {
	if m.shouldFail && m.failOn == distroName {
		return errors.New("export failed")
	}
	return nil
}

func TestBackupCommand(t *testing.T) {
	t.Run("command metadata", func(t *testing.T) {
		if RootCmd == nil {
			t.Fatal("RootCmd is nil")
		}

		backupCmd, _, err := RootCmd.Find([]string{"backup"})
		if err != nil {
			t.Fatalf("backup command not found: %v", err)
		}

		if backupCmd.Use != "backup <distro> [distro...]" {
			t.Errorf("expected Use='backup <distro> [distro...]', got '%s'", backupCmd.Use)
		}

		if backupCmd.Short == "" {
			t.Error("Short description is empty")
		}

		if backupCmd.Long == "" {
			t.Error("Long description is empty")
		}
	})

	t.Run("backup single distro success", func(t *testing.T) {
		skipIfNotWindows(t)
		mock := &mockBackuper{}
		out := new(bytes.Buffer)

		err := BackupDistrosCmd(context.Background(), mock, out, []string{"Ubuntu"}, "", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "Successfully backed up 1/1") {
			t.Errorf("expected success message in output, got:\n%s", output)
		}
	})

	t.Run("custom name with multiple distros returns error", func(t *testing.T) {
		mock := &mockBackuper{}
		out := new(bytes.Buffer)

		err := BackupDistrosCmd(context.Background(), mock, out, []string{"Ubuntu", "Debian"}, "custom", "")
		if err == nil {
			t.Fatal("expected error for custom name with multiple distros")
		}

		if !strings.Contains(err.Error(), "custom name can only be used when backing up a single distribution") {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("returns error when backup fails", func(t *testing.T) {
		skipIfNotWindows(t)
		mock := &mockBackuper{shouldFail: true, failOn: "Ubuntu"}
		out := new(bytes.Buffer)

		err := BackupDistrosCmd(context.Background(), mock, out, []string{"Ubuntu"}, "", "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "some backups failed") {
			t.Errorf("expected 'some backups failed' error, got: %v", err)
		}
	})

	t.Run("prints failure for failed backup", func(t *testing.T) {
		skipIfNotWindows(t)
		mock := &mockBackuper{shouldFail: true, failOn: "Ubuntu"}
		out := new(bytes.Buffer)

		BackupDistrosCmd(context.Background(), mock, out, []string{"Ubuntu"}, "", "")

		output := out.String()
		if !strings.Contains(output, "✗") {
			t.Errorf("expected failure indicator in output, got:\n%s", output)
		}
	})

	t.Run("flag defaults", func(t *testing.T) {
		backupCmd, _, err := RootCmd.Find([]string{"backup"})
		if err != nil {
			t.Fatalf("backup command not found: %v", err)
		}

		nameFlag := backupCmd.Flags().Lookup("name")
		if nameFlag == nil {
			t.Fatal("name flag not found")
		}

		backupDirFlag := backupCmd.Flags().Lookup("backup-dir")
		if backupDirFlag == nil {
			t.Fatal("backup-dir flag not found")
		}
	})
}
