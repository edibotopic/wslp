package wsl

import (
	"context"
	"fmt"

	gowsl "github.com/ubuntu/gowsl"
)

// UnregisterDistro unregisters a WSL distribution
func UnregisterDistro(ctx context.Context, name string) error {
	d := gowsl.NewDistro(ctx, name)
	
	registered, err := d.IsRegistered()
	if err != nil {
		return fmt.Errorf("failed to check if distro is registered: %w", err)
	}
	
	if !registered {
		return fmt.Errorf("distro %s is not registered", name)
	}
	
	if err := d.Unregister(); err != nil {
		return fmt.Errorf("failed to unregister distro: %w", err)
	}
	
	return nil
}
