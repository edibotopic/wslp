package cmd

import (
	"bytes"
	"context"
	"errors"
	"strings"
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

func TestDefaultCommand(t *testing.T) {
	t.Run("command metadata", func(t *testing.T) {
		if RootCmd == nil {
			t.Fatal("RootCmd is nil")
		}

		defaultCmd, _, err := RootCmd.Find([]string{"default"})
		if err != nil {
			t.Fatalf("default command not found: %v", err)
		}

		if defaultCmd.Use != "default" {
			t.Errorf("expected Use='default', got '%s'", defaultCmd.Use)
		}

		if defaultCmd.Short == "" {
			t.Error("Short description is empty")
		}

		if defaultCmd.Long == "" {
			t.Error("Long description is empty")
		}
	})

	t.Run("prints the default distro name", func(t *testing.T) {
		mock := &mockDefaultGetter{name: "Ubuntu"}
		out := new(bytes.Buffer)

		err := ShowDefault(context.Background(), mock, out)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		output := out.String()
		for _, phrase := range []string{"The default WSL distro is:", "Ubuntu"} {
			if !strings.Contains(output, phrase) {
				t.Errorf("expected output to contain %q, got:\n%s", phrase, output)
			}
		}
	})

	t.Run("returns error from getter", func(t *testing.T) {
		mock := &mockDefaultGetter{shouldFail: true}
		out := new(bytes.Buffer)

		err := ShowDefault(context.Background(), mock, out)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
