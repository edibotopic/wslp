package wsl

import (
	"context"
	"fmt"

	gowsl "github.com/ubuntu/gowsl"
)

// GetDefaultDistro retrieves the default WSL distribution
func GetDefaultDistro(ctx context.Context) (string, error) {
	defaultDistro, _, err := gowsl.DefaultDistro(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get default distro: %w", err)
	}

	return defaultDistro.Name(), nil
}
