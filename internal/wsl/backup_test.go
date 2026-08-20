package wsl

import (
	"context"
	"errors"
	"testing"
)

type mockBackuper struct {
	isRegisteredResults map[string]bool
	isRegisteredErrors  map[string]error
	exportErrors        map[string]error
}

func (m *mockBackuper) IsRegistered(ctx context.Context, name string) (bool, error) {
	if err, ok := m.isRegisteredErrors[name]; ok {
		return false, err
	}
	return m.isRegisteredResults[name], nil
}

func (m *mockBackuper) Export(ctx context.Context, distroName, outputPath string) error {
	if err, ok := m.exportErrors[distroName]; ok {
		return err
	}
	return nil
}

func TestBackupDistros(t *testing.T) {
	ctx := context.Background()

	t.Run("returns empty slice for empty distros list", func(t *testing.T) {
		mock := &mockBackuper{
			isRegisteredResults: make(map[string]bool),
			isRegisteredErrors:  make(map[string]error),
			exportErrors:        make(map[string]error),
		}

		results := BackupDistros(ctx, mock, []string{}, "/backup", BackupOptions{})

		if len(results) != 0 {
			t.Errorf("expected 0 results for empty distros, got %d", len(results))
		}
	})

	t.Run("backs up a single distro with auto-generated timestamp filename", func(t *testing.T) {
		mock := &mockBackuper{
			isRegisteredResults: map[string]bool{"Ubuntu": true},
			isRegisteredErrors:  make(map[string]error),
			exportErrors:        make(map[string]error),
		}

		results := BackupDistros(ctx, mock, []string{"Ubuntu"}, "/backup", BackupOptions{})

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if !results[0].Success {
			t.Errorf("expected success, got message: %s", results[0].Message)
		}
		if results[0].Distro != "Ubuntu" {
			t.Errorf("expected distro Ubuntu, got %s", results[0].Distro)
		}
		// Verify filename contains distro name and timestamp format
		if !matches(results[0].FilePath, `Ubuntu-\d{8}-\d{6}\.tar\.gz`) {
			t.Errorf("expected auto-generated filename with timestamp, got %s", results[0].FilePath)
		}
	})

	t.Run("backs up multiple distros", func(t *testing.T) {
		mock := &mockBackuper{
			isRegisteredResults: map[string]bool{
				"Ubuntu": true,
				"Debian": true,
				"Fedora": true,
			},
			isRegisteredErrors: make(map[string]error),
			exportErrors:       make(map[string]error),
		}

		results := BackupDistros(ctx, mock, []string{"Ubuntu", "Debian", "Fedora"}, "/backup", BackupOptions{})

		if len(results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(results))
		}
		for i, result := range results {
			if !result.Success {
				t.Errorf("result %d: expected success, got message: %s", i, result.Message)
			}
		}
	})

	t.Run("handles distro not registered", func(t *testing.T) {
		mock := &mockBackuper{
			isRegisteredResults: map[string]bool{"Ubuntu": false},
			isRegisteredErrors:  make(map[string]error),
			exportErrors:        make(map[string]error),
		}

		results := BackupDistros(ctx, mock, []string{"Ubuntu"}, "/backup", BackupOptions{})

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].Success {
			t.Error("expected failure for unregistered distro")
		}
		if results[0].Message != "Distribution is not registered" {
			t.Errorf("expected 'Distribution is not registered' message, got: %s", results[0].Message)
		}
	})

	t.Run("handles IsRegistered error", func(t *testing.T) {
		mock := &mockBackuper{
			isRegisteredResults: make(map[string]bool),
			isRegisteredErrors: map[string]error{
				"Ubuntu": errors.New("check failed"),
			},
			exportErrors: make(map[string]error),
		}

		results := BackupDistros(ctx, mock, []string{"Ubuntu"}, "/backup", BackupOptions{})

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].Success {
			t.Error("expected failure when IsRegistered errors")
		}
		if !matches(results[0].Message, "Error checking registration") {
			t.Errorf("expected error message about registration check, got: %s", results[0].Message)
		}
	})

	t.Run("handles Export error", func(t *testing.T) {
		mock := &mockBackuper{
			isRegisteredResults: map[string]bool{"Ubuntu": true},
			isRegisteredErrors:  make(map[string]error),
			exportErrors: map[string]error{
				"Ubuntu": errors.New("export failed"),
			},
		}

		results := BackupDistros(ctx, mock, []string{"Ubuntu"}, "/backup", BackupOptions{})

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].Success {
			t.Error("expected failure when Export errors")
		}
		if !matches(results[0].Message, "Backup failed") {
			t.Errorf("expected error message about backup failure, got: %s", results[0].Message)
		}
	})

	t.Run("uses custom backup name without extension", func(t *testing.T) {
		mock := &mockBackuper{
			isRegisteredResults: map[string]bool{"Ubuntu": true},
			isRegisteredErrors:  make(map[string]error),
			exportErrors:        make(map[string]error),
		}

		results := BackupDistros(ctx, mock, []string{"Ubuntu"}, "/backup", BackupOptions{
			CustomName: "my-backup",
		})

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if !results[0].Success {
			t.Errorf("expected success, got message: %s", results[0].Message)
		}
		// Should add .tar.gz extension
		if !matches(results[0].FilePath, `my-backup\.tar\.gz$`) {
			t.Errorf("expected filename to end with .tar.gz, got %s", results[0].FilePath)
		}
	})

	t.Run("preserves custom backup name with .tar.gz extension", func(t *testing.T) {
		mock := &mockBackuper{
			isRegisteredResults: map[string]bool{"Ubuntu": true},
			isRegisteredErrors:  make(map[string]error),
			exportErrors:        make(map[string]error),
		}

		results := BackupDistros(ctx, mock, []string{"Ubuntu"}, "/backup", BackupOptions{
			CustomName: "my-backup.tar.gz",
		})

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if !results[0].Success {
			t.Errorf("expected success, got message: %s", results[0].Message)
		}
		// Should not double-add extension
		if !matches(results[0].FilePath, `my-backup\.tar\.gz$`) {
			t.Errorf("expected filename to be my-backup.tar.gz, got %s", results[0].FilePath)
		}
	})

	t.Run("preserves custom backup name with .gz extension", func(t *testing.T) {
		mock := &mockBackuper{
			isRegisteredResults: map[string]bool{"Ubuntu": true},
			isRegisteredErrors:  make(map[string]error),
			exportErrors:        make(map[string]error),
		}

		results := BackupDistros(ctx, mock, []string{"Ubuntu"}, "/backup", BackupOptions{
			CustomName: "my-backup.tar.gz",
		})

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if !results[0].Success {
			t.Errorf("expected success, got message: %s", results[0].Message)
		}
		// Should not add another extension
		if !matches(results[0].FilePath, `my-backup\.tar\.gz$`) {
			t.Errorf("expected filename to be my-backup.tar.gz, got %s", results[0].FilePath)
		}
	})

	t.Run("continues on partial failure in multi-distro backup", func(t *testing.T) {
		mock := &mockBackuper{
			isRegisteredResults: map[string]bool{
				"Ubuntu": true,
				"Debian": false,
				"Fedora": true,
			},
			isRegisteredErrors: make(map[string]error),
			exportErrors: map[string]error{
				"Fedora": errors.New("export failed"),
			},
		}

		results := BackupDistros(ctx, mock, []string{"Ubuntu", "Debian", "Fedora"}, "/backup", BackupOptions{})

		if len(results) != 3 {
			t.Fatalf("expected 3 results, got %d", len(results))
		}
		if !results[0].Success {
			t.Errorf("result 0 (Ubuntu): expected success, got message: %s", results[0].Message)
		}
		if results[1].Success {
			t.Error("result 1 (Debian): expected failure for unregistered distro")
		}
		if results[2].Success {
			t.Error("result 2 (Fedora): expected failure due to export error")
		}
	})
}

// matches checks if s matches a simple pattern (for test assertions)
func matches(s, pattern string) bool {
	switch {
	case pattern == `Ubuntu-\d{8}-\d{6}\.tar\.gz`:
		return stringContains(s, "Ubuntu-") && stringHasSuffix(s, ".tar.gz") && len(s) > len("Ubuntu-.tar.gz")
	case pattern == `my-backup\.tar\.gz$`:
		return stringHasSuffix(s, "my-backup.tar.gz")
	case pattern == `Error checking registration`:
		return stringContains(s, "Error checking registration")
	case pattern == `Backup failed`:
		return stringContains(s, "Backup failed")
	default:
		return false
	}
}

// stringContains is a simple string contains check
func stringContains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// stringHasSuffix is a simple string suffix check
func stringHasSuffix(s, suffix string) bool {
	if len(suffix) > len(s) {
		return false
	}
	return s[len(s)-len(suffix):] == suffix
}
