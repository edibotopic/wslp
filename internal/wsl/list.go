package wsl

import (
	"context"
	"fmt"

	gowsl "github.com/ubuntu/gowsl"
)

// DistroInfo represents basic information about a WSL distro
type DistroInfo struct {
	Name string `json:"name"`
}

// ListDistros retrieves all registered WSL distributions
func ListDistros(ctx context.Context) ([]DistroInfo, error) {
	distroList, err := gowsl.RegisteredDistros(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get registered distros: %w", err)
	}

	result := make([]DistroInfo, len(distroList))
	for i := range distroList {
		result[i] = DistroInfo{
			Name: distroList[i].Name(),
		}
	}

	return result, nil
}
