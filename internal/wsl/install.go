package wsl

import (
	"context"
	"fmt"
	"io"
	"sync"

	gowsl "github.com/ubuntu/gowsl"
	"wslp/internal/config"
)

// InstallResult contains the result of a single distro installation
type InstallResult struct {
	Distro     string `json:"distro"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Registered bool   `json:"registered"`
}

// InstallDistros installs one or more WSL distributions
func InstallDistros(ctx context.Context, distros []string, concurrent bool) []InstallResult {
	if len(distros) == 0 {
		return []InstallResult{}
	}

	if !concurrent || len(distros) == 1 {
		results := make([]InstallResult, 0, len(distros))
		for _, distro := range distros {
			results = append(results, installOne(ctx, distro))
		}
		return results
	}

	results := make([]InstallResult, len(distros))
	sem := make(chan struct{}, config.GetMaxConcurrentInstalls())
	var wg sync.WaitGroup

	for i, distro := range distros {
		wg.Add(1)
		go func(idx int, name string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			results[idx] = installOne(ctx, name)
		}(i, distro)
	}

	wg.Wait()
	return results
}

// PrintInstallResults writes install results to out in a human-readable format
func PrintInstallResults(out io.Writer, results []InstallResult) {
	for _, r := range results {
		if !r.Success {
			fmt.Fprintf(out, "Error installing %s: %s\n", r.Distro, r.Message)
			continue
		}
		fmt.Fprintf(out, "Successfully installed %s\n", r.Distro)
		if r.Registered {
			fmt.Fprintf(out, "%s was downloaded in the modern format and is already registered\n", r.Distro)
			fmt.Fprintf(out, "Launch with wsl -d %s\n", r.Distro)
		} else {
			fmt.Fprintf(out, "%s was downloaded in the classic format and must be registered\n", r.Distro)
			fmt.Fprintf(out, "Register with wsl --register %s\n", r.Distro)
		}
	}
}


func installOne(ctx context.Context, distro string) InstallResult {
	result := InstallResult{Distro: distro}

	if err := gowsl.Install(ctx, distro); err != nil {
		result.Message = err.Error()
		return result
	}

	result.Success = true

	d := gowsl.NewDistro(ctx, distro)
	registered, err := d.IsRegistered()
	if err != nil {
		result.Message = fmt.Sprintf("Successfully installed (Error checking registration: %v)", err)
	} else {
		result.Registered = registered
		if registered {
			result.Message = "Successfully installed (modern format, already registered)"
		} else {
			result.Message = "Successfully installed (classic format, must be registered with: wsl --register " + distro + ")"
		}
	}

	return result
}
