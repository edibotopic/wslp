package wsl

import (
	"context"
	"errors"
	"testing"
)

type mockDefaultGetter struct {
	name       string
	shouldFail bool
}

func (m *mockDefaultGetter) GetDefault(ctx context.Context) (string, error) {
	if m.shouldFail {
		return "", errors.New("mock error")
	}
	return m.name, nil
}

type mockDefaultSetter struct {
	registered bool
	checkFails bool
	setFails   bool
	setName    string
}

func (m *mockDefaultSetter) IsRegistered(ctx context.Context, name string) (bool, error) {
	if m.checkFails {
		return false, errors.New("mock check error")
	}
	return m.registered, nil
}

func (m *mockDefaultSetter) SetAsDefault(ctx context.Context, name string) error {
	if m.setFails {
		return errors.New("mock set error")
	}
	m.setName = name
	return nil
}

func TestGetDefaultDistro(t *testing.T) {
	t.Run("returns the default distro name", func(t *testing.T) {
		mock := &mockDefaultGetter{name: "Ubuntu"}

		result, err := GetDefaultDistro(context.Background(), mock)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result != "Ubuntu" {
			t.Errorf("expected Ubuntu, got %s", result)
		}
	})

	t.Run("returns error from getter", func(t *testing.T) {
		mock := &mockDefaultGetter{shouldFail: true}

		_, err := GetDefaultDistro(context.Background(), mock)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestSetDefaultDistro(t *testing.T) {
	t.Run("sets a registered distro as default", func(t *testing.T) {
		mock := &mockDefaultSetter{registered: true}

		err := SetDefaultDistro(context.Background(), "Ubuntu", mock)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if mock.setName != "Ubuntu" {
			t.Errorf("expected Ubuntu to be set as default, got %s", mock.setName)
		}
	})

	t.Run("returns error if distro is not registered", func(t *testing.T) {
		mock := &mockDefaultSetter{registered: false}

		err := SetDefaultDistro(context.Background(), "Ubuntu", mock)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns error from IsRegistered", func(t *testing.T) {
		mock := &mockDefaultSetter{checkFails: true}

		err := SetDefaultDistro(context.Background(), "Ubuntu", mock)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns error from SetAsDefault", func(t *testing.T) {
		mock := &mockDefaultSetter{registered: true, setFails: true}

		err := SetDefaultDistro(context.Background(), "Ubuntu", mock)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
