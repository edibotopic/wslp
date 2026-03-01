package wsl

import (
	"context"
	"fmt"
	"os/exec"

	gowsl "github.com/ubuntu/gowsl"
)

// LaunchInteractive launches an interactive shell for the distro (blocking)
// This is suitable for CLI usage
func LaunchInteractive(ctx context.Context, distroName string) error {
	distro := gowsl.NewDistro(ctx, distroName)

	// Verify distro is registered
	registered, err := distro.IsRegistered()
	if err != nil {
		return fmt.Errorf("error checking registration: %w", err)
	}
	if !registered {
		return fmt.Errorf("distribution %s is not registered", distroName)
	}

	// Launch interactive shell (this is blocking)
	return distro.Shell()
}

// LaunchInTerminal launches the distro in a new terminal window (non-blocking)
// This is suitable for GUI usage
func LaunchInTerminal(ctx context.Context, distroName string) error {
	distro := gowsl.NewDistro(ctx, distroName)

	// Verify distro is registered
	registered, err := distro.IsRegistered()
	if err != nil {
		return fmt.Errorf("error checking registration: %w", err)
	}
	if !registered {
		return fmt.Errorf("distribution %s is not registered", distroName)
	}

	// Try to use Windows Terminal first
	if isWindowsTerminalAvailable() {
		// Use ~ to start in the distro's home directory
		cmd := exec.CommandContext(ctx, "wt.exe", "wsl.exe", "-d", distroName, "--cd", "~")
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to launch Windows Terminal: %w", err)
		}
		return nil
	}

	// Fallback to wsl.exe (opens in conhost)
	cmd := exec.CommandContext(ctx, "wsl.exe", "-d", distroName, "--cd", "~")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch wsl.exe: %w", err)
	}
	return nil
}

// isWindowsTerminalAvailable checks if Windows Terminal is available
func isWindowsTerminalAvailable() bool {
	_, err := exec.LookPath("wt.exe")
	return err == nil
}
