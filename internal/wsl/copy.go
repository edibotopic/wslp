package wsl

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	gowsl "github.com/ubuntu/gowsl"
)

// CopyResult contains the result of copying a distro
type CopyResult struct {
	Source  string `json:"source"`
	NewName string `json:"newName"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Copier interface for copying a distro
type Copier interface {
	IsRegistered(ctx context.Context, name string) (bool, error)
	Export(ctx context.Context, distroName, outputPath string) error
	Import(ctx context.Context, newName, tarPath, installDir string) error
}

// RealCopier implements Copier using gowsl and wsl.exe
type RealCopier struct{}

// IsRegistered checks if a distro is registered
func (r RealCopier) IsRegistered(ctx context.Context, name string) (bool, error) {
	return RealBackuper{}.IsRegistered(ctx, name)
}

// Export exports a distro to a tar.gz file
func (r RealCopier) Export(ctx context.Context, distroName, outputPath string) error {
	return RealBackuper{}.Export(ctx, distroName, outputPath)
}

// Import imports a distro from a tar file using gowsl
func (r RealCopier) Import(ctx context.Context, newName, tarPath, installDir string) error {
	_, err := gowsl.Import(ctx, newName, tarPath, installDir)
	return err
}

// CopyDistro copies a distro by exporting it to a temp file and importing it under a new name.
// installDir is where WSL stores the new distro's virtual disk; if empty it defaults to
// %USERPROFILE%\WSLCopies\<newName>.
func CopyDistro(ctx context.Context, c Copier, source, newName, installDir string) CopyResult {
	result := CopyResult{
		Source:  source,
		NewName: newName,
		Success: false,
	}

	if newName == "" {
		result.Message = "New name cannot be empty"
		return result
	}

	// Check source exists
	registered, err := c.IsRegistered(ctx, source)
	if err != nil {
		result.Message = fmt.Sprintf("Error checking source distro: %v", err)
		return result
	}
	if !registered {
		result.Message = fmt.Sprintf("Source distro %q is not registered", source)
		return result
	}

	// Check new name is not already taken
	exists, err := c.IsRegistered(ctx, newName)
	if err != nil {
		result.Message = fmt.Sprintf("Error checking new name: %v", err)
		return result
	}
	if exists {
		result.Message = fmt.Sprintf("Distro %q already exists", newName)
		return result
	}

	// Resolve install dir
	if installDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			result.Message = fmt.Sprintf("Failed to resolve home directory: %v", err)
			return result
		}
		installDir = filepath.Join(home, "WSLCopies", newName)
	}

	// Create install dir
	if err := os.MkdirAll(installDir, 0755); err != nil {
		result.Message = fmt.Sprintf("Failed to create install directory: %v", err)
		return result
	}

	// Export to a temp file
	timestamp := time.Now().Format("20060102-150405")
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("wslp-copy-%s-%s.tar.gz", source, timestamp))
	defer os.Remove(tmpFile)

	if err := c.Export(ctx, source, tmpFile); err != nil {
		result.Message = fmt.Sprintf("Export failed: %v", err)
		return result
	}

	// Import under the new name
	if err := c.Import(ctx, newName, tmpFile, installDir); err != nil {
		result.Message = fmt.Sprintf("Import failed: %v", err)
		return result
	}

	result.Success = true
	result.Message = fmt.Sprintf("Successfully copied %s to %s", source, newName)
	return result
}
