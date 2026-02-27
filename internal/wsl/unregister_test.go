package wsl

import (
	"context"
	"errors"
	"testing"
)

type mockUnregisterer struct {
	registered   bool
	checkFails   bool
	unregFails   bool
	unregistered []string
}

func (m *mockUnregisterer) IsRegistered(ctx context.Context, name string) (bool, error) {
	if m.checkFails {
		return false, errors.New("mock check error")
	}
	return m.registered, nil
}

func (m *mockUnregisterer) Unregister(ctx context.Context, name string) error {
	if m.unregFails {
		return errors.New("mock unregister error")
	}
	m.unregistered = append(m.unregistered, name)
	return nil
}

func TestUnregisterDistro(t *testing.T) {
	t.Run("unregisters a registered distro", func(t *testing.T) {
		mock := &mockUnregisterer{registered: true}

		err := UnregisterDistro(context.Background(), "Ubuntu", mock)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(mock.unregistered) != 1 || mock.unregistered[0] != "Ubuntu" {
			t.Errorf("expected Ubuntu to be unregistered, got: %v", mock.unregistered)
		}
	})

	t.Run("returns error if distro is not registered", func(t *testing.T) {
		mock := &mockUnregisterer{registered: false}

		err := UnregisterDistro(context.Background(), "Ubuntu", mock)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns error from IsRegistered", func(t *testing.T) {
		mock := &mockUnregisterer{checkFails: true}

		err := UnregisterDistro(context.Background(), "Ubuntu", mock)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns error from Unregister", func(t *testing.T) {
		mock := &mockUnregisterer{registered: true, unregFails: true}

		err := UnregisterDistro(context.Background(), "Ubuntu", mock)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
