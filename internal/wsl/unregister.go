package wsl

import (
	"context"
	"fmt"

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
	d := gowsl.NewDistro(ctx, name)
	return d.IsRegistered()
}

func (r RealUnregisterer) Unregister(ctx context.Context, name string) error {
	d := gowsl.NewDistro(ctx, name)
	return d.Unregister()
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
