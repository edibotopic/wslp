package wsl

import (
	"context"
	"fmt"

	gowsl "github.com/ubuntu/gowsl"
)

// InstallResult contains the result of a single distro installation
type InstallResult struct {
	Distro     string `json:"distro"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Registered bool   `json:"registered"`
}

// InstallDistros installs one or more WSL distributions
func InstallDistros(ctx context.Context, distros []string) []InstallResult {
	results := make([]InstallResult, 0, len(distros))

	if len(distros) == 0 {
		return results
	}

	// TODO: try to install distros concurrently
	for _, distro := range distros {
		result := InstallResult{
			Distro:  distro,
			Success: false,
		}

		err := gowsl.Install(ctx, distro)
		if err != nil {
			result.Message = fmt.Sprintf("Error installing: %v", err)
			results = append(results, result)
			continue
		}

		result.Success = true
		result.Message = "Successfully installed"

		d := gowsl.NewDistro(ctx, distro)
		registered, err := d.IsRegistered()
		if err != nil {
			result.Message += fmt.Sprintf(" (Error checking registration: %v)", err)
		} else {
			result.Registered = registered
			if registered {
				result.Message += " (modern format, already registered)"
			} else {
				result.Message += " (classic format, must be registered with: wsl --register " + distro + ")"
			}
		}

		results = append(results, result)
	}

	return results
}
