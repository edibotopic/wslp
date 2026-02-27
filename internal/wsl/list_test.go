package wsl

import (
	"context"
	"errors"
	"testing"
)

type mockLister struct {
	names      []string
	shouldFail bool
}

func (m *mockLister) List(ctx context.Context) ([]string, error) {
	if m.shouldFail {
		return nil, errors.New("mock error")
	}
	return m.names, nil
}

func TestListDistros(t *testing.T) {
	t.Run("returns registered distros", func(t *testing.T) {
		mock := &mockLister{names: []string{"Ubuntu", "Debian"}}

		result, err := ListDistros(context.Background(), mock)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 distros, got %d", len(result))
		}
		if result[0].Name != "Ubuntu" {
			t.Errorf("expected Ubuntu, got %s", result[0].Name)
		}
		if result[1].Name != "Debian" {
			t.Errorf("expected Debian, got %s", result[1].Name)
		}
	})

	t.Run("returns empty slice when no distros registered", func(t *testing.T) {
		mock := &mockLister{names: []string{}}

		result, err := ListDistros(context.Background(), mock)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 0 {
			t.Errorf("expected 0 distros, got %d", len(result))
		}
	})

	t.Run("propagates error from lister", func(t *testing.T) {
		mock := &mockLister{shouldFail: true}

		_, err := ListDistros(context.Background(), mock)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
