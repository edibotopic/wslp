package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

type mockUnregisterer struct {
	shouldFail bool
	failOn     string // specific distro to fail on
}

func (m *mockUnregisterer) IsRegistered(ctx context.Context, name string) (bool, error) {
	if m.shouldFail && m.failOn == name {
		return false, errors.New("mock error")
	}
	return true, nil
}

func (m *mockUnregisterer) Unregister(ctx context.Context, name string) error {
	if m.shouldFail && m.failOn == name {
		return errors.New("unregister failed")
	}
	return nil
}

func TestUnregisterCommand(t *testing.T) {
	t.Run("command metadata", func(t *testing.T) {
		if RootCmd == nil {
			t.Fatal("RootCmd is nil")
		}

		unregisterCmd, _, err := RootCmd.Find([]string{"unregister"})
		if err != nil {
			t.Fatalf("unregister command not found: %v", err)
		}

		if unregisterCmd.Use != "unregister <distro> [distro...]" {
			t.Errorf("expected Use='unregister <distro> [distro...]', got '%s'", unregisterCmd.Use)
		}

		if unregisterCmd.Short == "" {
			t.Error("Short description is empty")
		}

		if unregisterCmd.Long == "" {
			t.Error("Long description is empty")
		}

		// Check aliases
		hasDelete := false
		hasRemove := false
		for _, alias := range unregisterCmd.Aliases {
			if alias == "delete" {
				hasDelete = true
			}
			if alias == "remove" {
				hasRemove = true
			}
		}
		if !hasDelete {
			t.Error("expected 'delete' alias")
		}
		if !hasRemove {
			t.Error("expected 'remove' alias")
		}
	})

	t.Run("unregister single distro success", func(t *testing.T) {
		mock := &mockUnregisterer{}
		out := new(bytes.Buffer)

		err := UnregisterDistrosCmd(context.Background(), mock, out, []string{"Ubuntu"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "Successfully unregistered 1/1") {
			t.Errorf("expected success message in output, got:\n%s", output)
		}

		if !strings.Contains(output, "✓") {
			t.Errorf("expected success indicator in output, got:\n%s", output)
		}
	})

	t.Run("unregister multiple distros success", func(t *testing.T) {
		mock := &mockUnregisterer{}
		out := new(bytes.Buffer)

		err := UnregisterDistrosCmd(context.Background(), mock, out, []string{"Ubuntu", "Debian"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "Successfully unregistered 2/2") {
			t.Errorf("expected success message for 2 distros, got:\n%s", output)
		}
	})

	t.Run("returns error when unregister fails", func(t *testing.T) {
		mock := &mockUnregisterer{shouldFail: true, failOn: "Ubuntu"}
		out := new(bytes.Buffer)

		err := UnregisterDistrosCmd(context.Background(), mock, out, []string{"Ubuntu"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "some unregistrations failed") {
			t.Errorf("expected 'some unregistrations failed' error, got: %v", err)
		}
	})

	t.Run("prints failure for failed unregister", func(t *testing.T) {
		mock := &mockUnregisterer{shouldFail: true, failOn: "Ubuntu"}
		out := new(bytes.Buffer)

		UnregisterDistrosCmd(context.Background(), mock, out, []string{"Ubuntu"})

		output := out.String()
		if !strings.Contains(output, "✗") {
			t.Errorf("expected failure indicator in output, got:\n%s", output)
		}
	})

	t.Run("partial success", func(t *testing.T) {
		mock := &mockUnregisterer{shouldFail: true, failOn: "Debian"}
		out := new(bytes.Buffer)

		err := UnregisterDistrosCmd(context.Background(), mock, out, []string{"Ubuntu", "Debian"})
		if err == nil {
			t.Fatal("expected error for partial failure")
		}

		output := out.String()
		if !strings.Contains(output, "Successfully unregistered 1/2") {
			t.Errorf("expected partial success message, got:\n%s", output)
		}
	})
}
