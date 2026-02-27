package wsl

import (
	"context"
	"fmt"

	gowsl "github.com/ubuntu/gowsl"
)

// DefaultGetter retrieves the name of the default WSL distribution
type DefaultGetter interface {
	GetDefault(ctx context.Context) (string, error)
}

// DefaultSetter checks registration status and sets a WSL distribution as default
type DefaultSetter interface {
	IsRegistered(ctx context.Context, name string) (bool, error)
	SetAsDefault(ctx context.Context, name string) error
}

// RealDefaultGetter uses the actual gowsl library
type RealDefaultGetter struct{}

func (r RealDefaultGetter) GetDefault(ctx context.Context) (string, error) {
	d, _, err := gowsl.DefaultDistro(ctx)
	if err != nil {
		return "", err
	}
	return d.Name(), nil
}

// RealDefaultSetter uses the actual gowsl library
type RealDefaultSetter struct{}

func (r RealDefaultSetter) IsRegistered(ctx context.Context, name string) (bool, error) {
	d := gowsl.NewDistro(ctx, name)
	return d.IsRegistered()
}

func (r RealDefaultSetter) SetAsDefault(ctx context.Context, name string) error {
	d := gowsl.NewDistro(ctx, name)
	return d.SetAsDefault()
}

// GetDefaultDistro retrieves the default WSL distribution
func GetDefaultDistro(ctx context.Context, g DefaultGetter) (string, error) {
	name, err := g.GetDefault(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get default distro: %w", err)
	}
	return name, nil
}

// SetDefaultDistro sets a WSL distribution as the default
func SetDefaultDistro(ctx context.Context, name string, s DefaultSetter) error {
	registered, err := s.IsRegistered(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to check if distro is registered: %w", err)
	}

	if !registered {
		return fmt.Errorf("distro %s is not registered", name)
	}

	if err := s.SetAsDefault(ctx, name); err != nil {
		return fmt.Errorf("failed to set default distro: %w", err)
	}

	return nil
}
