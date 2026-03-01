package wsl

import (
	"context"
	"fmt"
	"strings"

	gowsl "github.com/ubuntu/gowsl"
)

// Unregisterer checks registration status and unregisters a WSL distribution
type Unregisterer interface {
	IsRegistered(ctx context.Context, name string) (bool, error)
	Unregister(ctx context.Context, name string) error
}

// RealUnregisterer uses the actual gowsl library
type RealUnregisterer struct{}

func (r RealUnregisterer) IsRegistered(ctx context.Context, name string) (bool, error) {
	distros, err := gowsl.RegisteredDistros(ctx)
	if err != nil {
		return false, err
	}
	// NOTE: Case-insensitive comparison. Distros with matching names but different
	// casings (e.g., "Ubuntu" and "ubuntu") is atypical but possible - this would
	// match the first one found.
	for _, d := range distros {
		if strings.EqualFold(d.Name(), name) {
			return true, nil
		}
	}
	return false, nil
}

func (r RealUnregisterer) Unregister(ctx context.Context, name string) error {
	d := gowsl.NewDistro(ctx, name)
	return d.Unregister()
}

// UnregisterResult contains the result of unregistering a distro
type UnregisterResult struct {
	Distro  string `json:"distro"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UnregisterDistro unregisters a WSL distribution
func UnregisterDistro(ctx context.Context, name string, u Unregisterer) error {
	registered, err := u.IsRegistered(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to check if distro is registered: %w", err)
	}

	if !registered {
		return fmt.Errorf("distro %s is not registered", name)
	}

	if err := u.Unregister(ctx, name); err != nil {
		return fmt.Errorf("failed to unregister distro: %w", err)
	}

	return nil
}

// UnregisterDistros unregisters one or more WSL distributions
func UnregisterDistros(ctx context.Context, u Unregisterer, distros []string) []UnregisterResult {
	results := make([]UnregisterResult, 0, len(distros))

	for _, distroName := range distros {
		result := UnregisterResult{
			Distro:  distroName,
			Success: false,
		}

		registered, err := u.IsRegistered(ctx, distroName)
		if err != nil {
			result.Message = fmt.Sprintf("Error checking registration: %v", err)
			results = append(results, result)
			continue
		}

		if !registered {
			result.Message = fmt.Sprintf("Distro %s is not registered", distroName)
			results = append(results, result)
			continue
		}

		if err := u.Unregister(ctx, distroName); err != nil {
			result.Message = fmt.Sprintf("Failed to unregister: %v", err)
			results = append(results, result)
			continue
		}

		result.Success = true
		result.Message = fmt.Sprintf("Successfully unregistered %s", distroName)
		results = append(results, result)
	}

	return results
}
