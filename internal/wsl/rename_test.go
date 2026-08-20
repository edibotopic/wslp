package wsl

import (
	"context"
	"errors"
	"testing"
)

type mockRenamer struct {
	isRegisteredResults map[string]bool
	isRegisteredErrors  map[string]error
	getDistroGUIDError  error
	renameInRegistryErr error
	renamedName         string
}

func (m *mockRenamer) IsRegistered(ctx context.Context, name string) (bool, error) {
	if err, ok := m.isRegisteredErrors[name]; ok {
		return false, err
	}
	if result, ok := m.isRegisteredResults[name]; ok {
		return result, nil
	}
	return false, nil
}

func (m *mockRenamer) GetDistroGUID(ctx context.Context, name string) (string, error) {
	if m.getDistroGUIDError != nil {
		return "", m.getDistroGUIDError
	}
	return "{12345678-1234-1234-1234-123456789012}", nil
}

func (m *mockRenamer) RenameInRegistry(guid, newName string) error {
	if m.renameInRegistryErr != nil {
		return m.renameInRegistryErr
	}
	m.renamedName = newName
	return nil
}

func TestRenameDistro(t *testing.T) {
	t.Run("returns error if new name is empty", func(t *testing.T) {
		mock := &mockRenamer{}

		result := RenameDistro(context.Background(), mock, "Ubuntu", "")

		if result.Success {
			t.Errorf("expected failure, got success")
		}
		if result.Message == "" {
			t.Errorf("expected error message, got empty")
		}
	})

	t.Run("returns error if old distro is not registered", func(t *testing.T) {
		mock := &mockRenamer{
			isRegisteredResults: map[string]bool{
				"Ubuntu": false,
			},
		}

		result := RenameDistro(context.Background(), mock, "Ubuntu", "NewName")

		if result.Success {
			t.Errorf("expected failure, got success")
		}
		if result.Message == "" {
			t.Errorf("expected error message, got empty")
		}
	})

	t.Run("returns error from IsRegistered when checking old name", func(t *testing.T) {
		mock := &mockRenamer{
			isRegisteredErrors: map[string]error{
				"Ubuntu": errors.New("registry error"),
			},
		}

		result := RenameDistro(context.Background(), mock, "Ubuntu", "NewName")

		if result.Success {
			t.Errorf("expected failure, got success")
		}
	})

	t.Run("returns error if new name already exists", func(t *testing.T) {
		mock := &mockRenamer{
			isRegisteredResults: map[string]bool{
				"Ubuntu":  true,
				"NewName": true,
			},
		}

		result := RenameDistro(context.Background(), mock, "Ubuntu", "NewName")

		if result.Success {
			t.Errorf("expected failure, got success")
		}
		if result.Message == "" {
			t.Errorf("expected error message, got empty")
		}
	})

	t.Run("returns error from IsRegistered when checking new name", func(t *testing.T) {
		mock := &mockRenamer{
			isRegisteredResults: map[string]bool{
				"Ubuntu": true,
			},
			isRegisteredErrors: map[string]error{
				"NewName": errors.New("registry error"),
			},
		}

		result := RenameDistro(context.Background(), mock, "Ubuntu", "NewName")

		if result.Success {
			t.Errorf("expected failure, got success")
		}
	})

	t.Run("returns error from GetDistroGUID", func(t *testing.T) {
		mock := &mockRenamer{
			isRegisteredResults: map[string]bool{
				"Ubuntu":  true,
				"NewName": false,
			},
			getDistroGUIDError: errors.New("failed to get GUID"),
		}

		result := RenameDistro(context.Background(), mock, "Ubuntu", "NewName")

		if result.Success {
			t.Errorf("expected failure, got success")
		}
	})

	t.Run("returns error from RenameInRegistry", func(t *testing.T) {
		mock := &mockRenamer{
			isRegisteredResults: map[string]bool{
				"Ubuntu":  true,
				"NewName": false,
			},
			renameInRegistryErr: errors.New("registry write failed"),
		}

		result := RenameDistro(context.Background(), mock, "Ubuntu", "NewName")

		if result.Success {
			t.Errorf("expected failure, got success")
		}
	})

	t.Run("successfully renames a distro", func(t *testing.T) {
		mock := &mockRenamer{
			isRegisteredResults: map[string]bool{
				"Ubuntu":  true,
				"NewName": false,
			},
		}

		result := RenameDistro(context.Background(), mock, "Ubuntu", "NewName")

		if !result.Success {
			t.Errorf("expected success, got failure: %s", result.Message)
		}
		if result.OldName != "Ubuntu" {
			t.Errorf("expected OldName to be Ubuntu, got %s", result.OldName)
		}
		if result.NewName != "NewName" {
			t.Errorf("expected NewName to be NewName, got %s", result.NewName)
		}
		if mock.renamedName != "NewName" {
			t.Errorf("expected distro to be renamed to NewName, got %s", mock.renamedName)
		}
	})
}
