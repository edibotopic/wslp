package wsl

import (
	"context"
	"errors"
	"testing"
)

type mockTerminator struct {
	isRegisteredResults map[string]bool
	isRegisteredErrors  map[string]error
	terminateErrors     map[string]error
	terminatedDistros   []string
}

func (m *mockTerminator) IsRegistered(ctx context.Context, name string) (bool, error) {
	if err, ok := m.isRegisteredErrors[name]; ok {
		return false, err
	}
	if result, ok := m.isRegisteredResults[name]; ok {
		return result, nil
	}
	return false, nil
}

func (m *mockTerminator) Terminate(ctx context.Context, name string) error {
	if err, ok := m.terminateErrors[name]; ok {
		return err
	}
	m.terminatedDistros = append(m.terminatedDistros, name)
	return nil
}

func TestTerminateDistros(t *testing.T) {
	t.Run("returns empty slice for empty distros list", func(t *testing.T) {
		mock := &mockTerminator{}

		results := TerminateDistros(context.Background(), mock, []string{})

		if len(results) != 0 {
			t.Errorf("expected empty results, got %d results", len(results))
		}
	})

	t.Run("successfully terminates a single distro", func(t *testing.T) {
		mock := &mockTerminator{
			isRegisteredResults: map[string]bool{
				"Ubuntu": true,
			},
		}

		results := TerminateDistros(context.Background(), mock, []string{"Ubuntu"})

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
		if !results[0].Success {
			t.Errorf("expected success, got failure: %s", results[0].Message)
		}
		if results[0].Distro != "Ubuntu" {
			t.Errorf("expected distro to be Ubuntu, got %s", results[0].Distro)
		}
		if len(mock.terminatedDistros) != 1 || mock.terminatedDistros[0] != "Ubuntu" {
			t.Errorf("expected Ubuntu to be terminated, got %v", mock.terminatedDistros)
		}
	})

	t.Run("successfully terminates multiple distros", func(t *testing.T) {
		mock := &mockTerminator{
			isRegisteredResults: map[string]bool{
				"Ubuntu": true,
				"Debian": true,
				"Fedora": true,
			},
		}

		results := TerminateDistros(context.Background(), mock, []string{"Ubuntu", "Debian", "Fedora"})

		if len(results) != 3 {
			t.Errorf("expected 3 results, got %d", len(results))
		}
		for i, result := range results {
			if !result.Success {
				t.Errorf("result[%d] expected success, got failure: %s", i, result.Message)
			}
		}
		if len(mock.terminatedDistros) != 3 {
			t.Errorf("expected 3 distros terminated, got %d", len(mock.terminatedDistros))
		}
	})

	t.Run("handles mixed success and failure across multiple distros", func(t *testing.T) {
		mock := &mockTerminator{
			isRegisteredResults: map[string]bool{
				"Ubuntu": true,
				"Debian": true,
			},
			terminateErrors: map[string]error{
				"Debian": errors.New("terminate failed"),
			},
		}

		results := TerminateDistros(context.Background(), mock, []string{"Ubuntu", "Debian"})

		if len(results) != 2 {
			t.Errorf("expected 2 results, got %d", len(results))
		}
		if !results[0].Success {
			t.Errorf("expected Ubuntu to succeed, got failure: %s", results[0].Message)
		}
		if results[1].Success {
			t.Errorf("expected Debian to fail, got success")
		}
		if len(mock.terminatedDistros) != 1 || mock.terminatedDistros[0] != "Ubuntu" {
			t.Errorf("expected only Ubuntu to be terminated, got %v", mock.terminatedDistros)
		}
	})

	t.Run("returns error if distro is not registered", func(t *testing.T) {
		mock := &mockTerminator{
			isRegisteredResults: map[string]bool{
				"Ubuntu": false,
			},
		}

		results := TerminateDistros(context.Background(), mock, []string{"Ubuntu"})

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
		if results[0].Success {
			t.Errorf("expected failure for unregistered distro")
		}
		if len(mock.terminatedDistros) != 0 {
			t.Errorf("expected no distros to be terminated, got %v", mock.terminatedDistros)
		}
	})

	t.Run("returns error from IsRegistered", func(t *testing.T) {
		mock := &mockTerminator{
			isRegisteredErrors: map[string]error{
				"Ubuntu": errors.New("registry error"),
			},
		}

		results := TerminateDistros(context.Background(), mock, []string{"Ubuntu"})

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
		if results[0].Success {
			t.Errorf("expected failure, got success")
		}
		if len(mock.terminatedDistros) != 0 {
			t.Errorf("expected no distros to be terminated, got %v", mock.terminatedDistros)
		}
	})

	t.Run("returns error from Terminate", func(t *testing.T) {
		mock := &mockTerminator{
			isRegisteredResults: map[string]bool{
				"Ubuntu": true,
			},
			terminateErrors: map[string]error{
				"Ubuntu": errors.New("terminate error"),
			},
		}

		results := TerminateDistros(context.Background(), mock, []string{"Ubuntu"})

		if len(results) != 1 {
			t.Errorf("expected 1 result, got %d", len(results))
		}
		if results[0].Success {
			t.Errorf("expected failure, got success")
		}
		if len(mock.terminatedDistros) != 0 {
			t.Errorf("expected no distros to be terminated, got %v", mock.terminatedDistros)
		}
	})
}
