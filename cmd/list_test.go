package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
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

func TestListCommand(t *testing.T) {
	t.Run("command metadata", func(t *testing.T) {
		if RootCmd == nil {
			t.Fatal("RootCmd is nil")
		}

		listCmd, _, err := RootCmd.Find([]string{"list"})
		if err != nil {
			t.Fatalf("list command not found: %v", err)
		}

		if listCmd.Use != "list" {
			t.Errorf("expected Use='list', got '%s'", listCmd.Use)
		}

		if listCmd.Short == "" {
			t.Error("Short description is empty")
		}

		if listCmd.Long == "" {
			t.Error("Long description is empty")
		}
	})

	t.Run("prints distros and count", func(t *testing.T) {
		mock := &mockLister{names: []string{"Ubuntu", "Debian"}}
		out := new(bytes.Buffer)

		err := ListDistros(context.Background(), mock, out)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()
		for _, phrase := range []string{"Finding registered distros", "2 distros are registered", "Ubuntu", "Debian"} {
			if !strings.Contains(output, phrase) {
				t.Errorf("expected output to contain %q, got:\n%s", phrase, output)
			}
		}
	})

	t.Run("prints zero count when no distros", func(t *testing.T) {
		mock := &mockLister{names: []string{}}
		out := new(bytes.Buffer)

		err := ListDistros(context.Background(), mock, out)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(out.String(), "0 distros are registered") {
			t.Errorf("expected output to mention 0 distros, got:\n%s", out.String())
		}
	})

	t.Run("returns error from lister", func(t *testing.T) {
		mock := &mockLister{shouldFail: true}
		out := new(bytes.Buffer)

		err := ListDistros(context.Background(), mock, out)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
