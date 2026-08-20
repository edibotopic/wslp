package wsl

import (
	"context"
	"errors"
	"testing"
)

type mockWorkshopRunner struct {
	output []byte
	err    error
}

func (m *mockWorkshopRunner) ListWorkshops(ctx context.Context, distro string) ([]byte, error) {
	return m.output, m.err
}

func TestGetWorkshops(t *testing.T) {
	t.Run("parses a well-formed table", func(t *testing.T) {
		mock := &mockWorkshopRunner{
			output: []byte(
				"/home/user/proj1  dev    Ready    -\n" +
					"/home/user/proj2  api    Stopped  -\n",
			),
		}

		result := GetWorkshops(context.Background(), "Ubuntu", mock)

		if len(result) != 2 {
			t.Fatalf("expected 2 workshops, got %d", len(result))
		}
		if result[0].Name != "dev" || result[0].Status != "Ready" {
			t.Errorf("unexpected first workshop: %+v", result[0])
		}
		if result[1].Name != "api" || result[1].Status != "Stopped" {
			t.Errorf("unexpected second workshop: %+v", result[1])
		}
	})

	t.Run("returns empty slice, not error, when workshop CLI is unavailable", func(t *testing.T) {
		mock := &mockWorkshopRunner{err: errors.New("exec: \"workshop\": executable file not found")}

		result := GetWorkshops(context.Background(), "Debian", mock)

		if result == nil {
			t.Fatal("expected non-nil empty slice")
		}
		if len(result) != 0 {
			t.Errorf("expected 0 workshops, got %d", len(result))
		}
	})

	t.Run("returns empty slice for blank output", func(t *testing.T) {
		mock := &mockWorkshopRunner{output: []byte("")}

		result := GetWorkshops(context.Background(), "Ubuntu", mock)

		if len(result) != 0 {
			t.Errorf("expected 0 workshops, got %d", len(result))
		}
	})

	t.Run("skips malformed lines", func(t *testing.T) {
		mock := &mockWorkshopRunner{output: []byte("garbage\ndev proj Ready notes here\n")}

		result := GetWorkshops(context.Background(), "Ubuntu", mock)

		if len(result) != 1 {
			t.Fatalf("expected 1 workshop, got %d", len(result))
		}
	})
}
