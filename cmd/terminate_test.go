package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

type mockTerminator struct {
	shouldFail bool
	failOn     string // specific distro to fail on
}

func (m *mockTerminator) IsRegistered(ctx context.Context, name string) (bool, error) {
	if m.shouldFail && m.failOn == name {
		return false, errors.New("mock error")
	}
	return true, nil
}

func (m *mockTerminator) Terminate(ctx context.Context, name string) error {
	if m.shouldFail && m.failOn == name {
		return errors.New("terminate failed")
	}
	return nil
}

func TestTerminateCommand(t *testing.T) {
	t.Run("command metadata", func(t *testing.T) {
		if RootCmd == nil {
			t.Fatal("RootCmd is nil")
		}

		terminateCmd, _, err := RootCmd.Find([]string{"terminate"})
		if err != nil {
			t.Fatalf("terminate command not found: %v", err)
		}

		if terminateCmd.Use != "terminate <distro> [distro...]" {
			t.Errorf("expected Use='terminate <distro> [distro...]', got '%s'", terminateCmd.Use)
		}

		if terminateCmd.Short == "" {
			t.Error("Short description is empty")
		}

		if terminateCmd.Long == "" {
			t.Error("Long description is empty")
		}

		// Check aliases
		hasStop := false
		hasKill := false
		for _, alias := range terminateCmd.Aliases {
			if alias == "stop" {
				hasStop = true
			}
			if alias == "kill" {
				hasKill = true
			}
		}
		if !hasStop {
			t.Error("expected 'stop' alias")
		}
		if !hasKill {
			t.Error("expected 'kill' alias")
		}
	})

	t.Run("terminate single distro success", func(t *testing.T) {
		mock := &mockTerminator{}
		out := new(bytes.Buffer)

		err := TerminateDistrosCmd(context.Background(), mock, out, []string{"Ubuntu"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "Successfully terminated 1/1") {
			t.Errorf("expected success message in output, got:\n%s", output)
		}

		if !strings.Contains(output, "✓") {
			t.Errorf("expected success indicator in output, got:\n%s", output)
		}
	})

	t.Run("terminate multiple distros success", func(t *testing.T) {
		mock := &mockTerminator{}
		out := new(bytes.Buffer)

		err := TerminateDistrosCmd(context.Background(), mock, out, []string{"Ubuntu", "Debian"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()
		if !strings.Contains(output, "Successfully terminated 2/2") {
			t.Errorf("expected success message for 2 distros, got:\n%s", output)
		}
	})

	t.Run("returns error when termination fails", func(t *testing.T) {
		mock := &mockTerminator{shouldFail: true, failOn: "Ubuntu"}
		out := new(bytes.Buffer)

		err := TerminateDistrosCmd(context.Background(), mock, out, []string{"Ubuntu"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "some terminations failed") {
			t.Errorf("expected 'some terminations failed' error, got: %v", err)
		}
	})

	t.Run("prints failure for failed termination", func(t *testing.T) {
		mock := &mockTerminator{shouldFail: true, failOn: "Ubuntu"}
		out := new(bytes.Buffer)

		TerminateDistrosCmd(context.Background(), mock, out, []string{"Ubuntu"})

		output := out.String()
		if !strings.Contains(output, "✗") {
			t.Errorf("expected failure indicator in output, got:\n%s", output)
		}
	})

	t.Run("partial success", func(t *testing.T) {
		mock := &mockTerminator{shouldFail: true, failOn: "Debian"}
		out := new(bytes.Buffer)

		err := TerminateDistrosCmd(context.Background(), mock, out, []string{"Ubuntu", "Debian"})
		if err == nil {
			t.Fatal("expected error for partial failure")
		}

		output := out.String()
		if !strings.Contains(output, "Successfully terminated 1/2") {
			t.Errorf("expected partial success message, got:\n%s", output)
		}
	})
}
