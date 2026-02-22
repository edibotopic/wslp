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

// SetDefaultDistro sets a WSL distribution as the default
func SetDefaultDistro(ctx context.Context, name string) error {
	d := gowsl.NewDistro(ctx, name)
	
	registered, err := d.IsRegistered()
	if err != nil {
		return fmt.Errorf("failed to check if distro is registered: %w", err)
	}
	
	if !registered {
		return fmt.Errorf("distro %s is not registered", name)
	}
	
	if err := d.SetAsDefault(); err != nil {
		return fmt.Errorf("failed to set default distro: %w", err)
	}
	
	return nil
}
