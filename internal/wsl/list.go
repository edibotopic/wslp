package wsl

import (
	"context"
	"fmt"

	gowsl "github.com/ubuntu/gowsl"
)

// DistroInfo represents basic information about a WSL distro
type DistroInfo struct {
	Name    string `json:"name"`
	State   string `json:"state"`
	Running bool   `json:"running"`
}

// Lister retrieves the names of registered WSL distributions
type Lister interface {
	List(ctx context.Context) ([]string, error)
}

// RealLister uses the actual gowsl library
type RealLister struct{}

func (r RealLister) List(ctx context.Context) ([]string, error) {
	distroList, err := gowsl.RegisteredDistros(ctx)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(distroList))
	for i := range distroList {
		names[i] = distroList[i].Name()
	}

	return names, nil
}

// ListDistros retrieves all registered WSL distributions with their state
func ListDistros(ctx context.Context, l Lister) ([]DistroInfo, error) {
	names, err := l.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get registered distros: %w", err)
	}

	result := make([]DistroInfo, len(names))
	for i, name := range names {
		distro := gowsl.NewDistro(ctx, name)
		state, err := distro.State()

		stateStr := "Unknown"
		running := false
		if err == nil {
			stateStr = state.String()
			running = (state != gowsl.Stopped)
		}

		result[i] = DistroInfo{
			Name:    name,
			State:   stateStr,
			Running: running,
		}
	}

	return result, nil
}
