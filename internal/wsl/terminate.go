package wsl

import (
	"context"
	"fmt"

	gowsl "github.com/ubuntu/gowsl"
)

// TerminateResult contains the result of terminating a distro
type TerminateResult struct {
	Distro  string `json:"distro"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Terminator interface for terminating distros
type Terminator interface {
	IsRegistered(ctx context.Context, name string) (bool, error)
	Terminate(ctx context.Context, name string) error
}

// RealTerminator implements Terminator using gowsl
type RealTerminator struct{}

// IsRegistered checks if a distro is registered
func (r RealTerminator) IsRegistered(ctx context.Context, name string) (bool, error) {
	distros, err := gowsl.RegisteredDistros(ctx)
	if err != nil {
		return false, err
	}
	// NOTE: Exact match comparison. Distros with matching names but different
	// casings (e.g., "Ubuntu" and "ubuntu") is atypical but possible.
	for _, d := range distros {
		if d.Name() == name {
			return true, nil
		}
	}
	return false, nil
}

// Terminate terminates a running distro
func (r RealTerminator) Terminate(ctx context.Context, name string) error {
	distro := gowsl.NewDistro(ctx, name)
	return distro.Terminate()
}

// TerminateDistros terminates one or more distros
func TerminateDistros(ctx context.Context, t Terminator, distros []string) []TerminateResult {
	results := make([]TerminateResult, 0, len(distros))

	if len(distros) == 0 {
		return results
	}

	for _, distroName := range distros {
		result := TerminateResult{
			Distro:  distroName,
			Success: false,
		}

		// Verify distro is registered
		registered, err := t.IsRegistered(ctx, distroName)
		if err != nil {
			result.Message = fmt.Sprintf("Error checking registration: %v", err)
			results = append(results, result)
			continue
		}
		if !registered {
			result.Message = "Distribution is not registered"
			results = append(results, result)
			continue
		}

		// Terminate the distro
		err = t.Terminate(ctx, distroName)
		if err != nil {
			result.Message = fmt.Sprintf("Failed to terminate: %v", err)
			results = append(results, result)
			continue
		}

		result.Success = true
		result.Message = "Successfully terminated"
		results = append(results, result)
	}

	return results
}
