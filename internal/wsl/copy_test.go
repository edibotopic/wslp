package wsl

import (
	"context"
	"errors"
	"testing"
)

type mockCopier struct {
	isRegisteredResults map[string]bool
	isRegisteredErrors  map[string]error
	exportErrors        map[string]error
	importErrors        map[string]error
}

func (m *mockCopier) IsRegistered(ctx context.Context, name string) (bool, error) {
	if err, ok := m.isRegisteredErrors[name]; ok {
		return false, err
	}
	return m.isRegisteredResults[name], nil
}

func (m *mockCopier) Export(ctx context.Context, distroName, outputPath string) error {
	if err, ok := m.exportErrors[distroName]; ok {
		return err
	}
	return nil
}

func (m *mockCopier) Import(ctx context.Context, newName, tarPath, installDir string) error {
	if err, ok := m.importErrors[newName]; ok {
		return err
	}
	return nil
}

func TestCopyDistro(t *testing.T) {
	ctx := context.Background()

	t.Run("returns error when new name is empty", func(t *testing.T) {
		mock := &mockCopier{
			isRegisteredResults: make(map[string]bool),
			isRegisteredErrors:  make(map[string]error),
			exportErrors:        make(map[string]error),
			importErrors:        make(map[string]error),
		}

		result := CopyDistro(ctx, mock, "Ubuntu", "", "/tmp/install")

		if result.Success {
			t.Error("expected failure for empty new name")
		}
		if result.Message != "New name cannot be empty" {
			t.Errorf("expected 'New name cannot be empty' message, got: %s", result.Message)
		}
	})

	t.Run("returns error when source distro is not registered", func(t *testing.T) {
		mock := &mockCopier{
			isRegisteredResults: map[string]bool{"Ubuntu": false},
			isRegisteredErrors:  make(map[string]error),
			exportErrors:        make(map[string]error),
			importErrors:        make(map[string]error),
		}

		result := CopyDistro(ctx, mock, "Ubuntu", "Ubuntu-Copy", "/tmp/install")

		if result.Success {
			t.Error("expected failure for unregistered source distro")
		}
		if !stringContains(result.Message, "not registered") {
			t.Errorf("expected message about source distro not registered, got: %s", result.Message)
		}
	})

	t.Run("returns error when new name already exists", func(t *testing.T) {
		mock := &mockCopier{
			isRegisteredResults: map[string]bool{
				"Ubuntu":      true,
				"Ubuntu-Copy": true,
			},
			isRegisteredErrors: make(map[string]error),
			exportErrors:       make(map[string]error),
			importErrors:       make(map[string]error),
		}

		result := CopyDistro(ctx, mock, "Ubuntu", "Ubuntu-Copy", "/tmp/install")

		if result.Success {
			t.Error("expected failure when new name already exists")
		}
		if !stringContains(result.Message, "already exists") {
			t.Errorf("expected message about new name already existing, got: %s", result.Message)
		}
	})

	t.Run("returns error when checking source registration fails", func(t *testing.T) {
		mock := &mockCopier{
			isRegisteredResults: make(map[string]bool),
			isRegisteredErrors: map[string]error{
				"Ubuntu": errors.New("check failed"),
			},
			exportErrors: make(map[string]error),
			importErrors: make(map[string]error),
		}

		result := CopyDistro(ctx, mock, "Ubuntu", "Ubuntu-Copy", "/tmp/install")

		if result.Success {
			t.Error("expected failure when source check errors")
		}
		if !stringContains(result.Message, "Error checking source distro") {
			t.Errorf("expected error message about source distro check, got: %s", result.Message)
		}
	})

	t.Run("returns error when checking new name registration fails", func(t *testing.T) {
		mock := &mockCopier{
			isRegisteredResults: map[string]bool{"Ubuntu": true},
			isRegisteredErrors: map[string]error{
				"Ubuntu-Copy": errors.New("check failed"),
			},
			exportErrors: make(map[string]error),
			importErrors: make(map[string]error),
		}

		result := CopyDistro(ctx, mock, "Ubuntu", "Ubuntu-Copy", "/tmp/install")

		if result.Success {
			t.Error("expected failure when new name check errors")
		}
		if !stringContains(result.Message, "Error checking new name") {
			t.Errorf("expected error message about new name check, got: %s", result.Message)
		}
	})

	t.Run("returns error when Export fails", func(t *testing.T) {
		mock := &mockCopier{
			isRegisteredResults: map[string]bool{
				"Ubuntu":      true,
				"Ubuntu-Copy": false,
			},
			isRegisteredErrors: make(map[string]error),
			exportErrors: map[string]error{
				"Ubuntu": errors.New("export failed"),
			},
			importErrors: make(map[string]error),
		}

		result := CopyDistro(ctx, mock, "Ubuntu", "Ubuntu-Copy", "/tmp/install")

		if result.Success {
			t.Error("expected failure when Export errors")
		}
		if !stringContains(result.Message, "Export failed") {
			t.Errorf("expected error message about export failure, got: %s", result.Message)
		}
	})

	t.Run("returns error when Import fails", func(t *testing.T) {
		mock := &mockCopier{
			isRegisteredResults: map[string]bool{
				"Ubuntu":      true,
				"Ubuntu-Copy": false,
			},
			isRegisteredErrors: make(map[string]error),
			exportErrors:       make(map[string]error),
			importErrors: map[string]error{
				"Ubuntu-Copy": errors.New("import failed"),
			},
		}

		result := CopyDistro(ctx, mock, "Ubuntu", "Ubuntu-Copy", "/tmp/install")

		if result.Success {
			t.Error("expected failure when Import errors")
		}
		if !stringContains(result.Message, "Import failed") {
			t.Errorf("expected error message about import failure, got: %s", result.Message)
		}
	})

	t.Run("successfully copies a distro with custom install directory", func(t *testing.T) {
		mock := &mockCopier{
			isRegisteredResults: map[string]bool{
				"Ubuntu":      true,
				"Ubuntu-Copy": false,
			},
			isRegisteredErrors: make(map[string]error),
			exportErrors:       make(map[string]error),
			importErrors:       make(map[string]error),
		}

		result := CopyDistro(ctx, mock, "Ubuntu", "Ubuntu-Copy", "/tmp/install")

		if !result.Success {
			t.Errorf("expected success, got message: %s", result.Message)
		}
		if result.Source != "Ubuntu" {
			t.Errorf("expected source Ubuntu, got %s", result.Source)
		}
		if result.NewName != "Ubuntu-Copy" {
			t.Errorf("expected new name Ubuntu-Copy, got %s", result.NewName)
		}
	})

	t.Run("successfully copies a distro with temp directory for install (using t.TempDir)", func(t *testing.T) {
		mock := &mockCopier{
			isRegisteredResults: map[string]bool{
				"Ubuntu":      true,
				"Ubuntu-Copy": false,
			},
			isRegisteredErrors: make(map[string]error),
			exportErrors:       make(map[string]error),
			importErrors:       make(map[string]error),
		}

		// Use a temp directory so the test is hermetic and self-cleaning
		tempDir := t.TempDir()

		result := CopyDistro(ctx, mock, "Ubuntu", "Ubuntu-Copy", tempDir)

		if !result.Success {
			t.Errorf("expected success, got message: %s", result.Message)
		}
		if result.Source != "Ubuntu" {
			t.Errorf("expected source Ubuntu, got %s", result.Source)
		}
		if result.NewName != "Ubuntu-Copy" {
			t.Errorf("expected new name Ubuntu-Copy, got %s", result.NewName)
		}
	})
}
