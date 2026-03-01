package wsl

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	gowsl "github.com/ubuntu/gowsl"
)

// BackupResult contains the result of a single distro backup
type BackupResult struct {
	Distro   string `json:"distro"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	FilePath string `json:"filePath"`
}

// BackupOptions contains options for backup operations
type BackupOptions struct {
	// CustomName is an optional custom name for the backup file
	// If empty, auto-generates name with timestamp
	CustomName string
}

// Backuper interface for backing up distros
type Backuper interface {
	IsRegistered(ctx context.Context, name string) (bool, error)
	Export(ctx context.Context, distroName, outputPath string) error
}

// RealBackuper implements Backuper using wsl.exe
type RealBackuper struct{}

// IsRegistered checks if a distro is registered
func (r RealBackuper) IsRegistered(ctx context.Context, name string) (bool, error) {
	distros, err := gowsl.RegisteredDistros(ctx)
	if err != nil {
		return false, err
	}
	for _, d := range distros {
		if d.Name() == name {
			return true, nil
		}
	}
	return false, nil
}

// Export backs up a distro using wsl.exe --export with tar.gz format
func (r RealBackuper) Export(ctx context.Context, distroName, outputPath string) error {
	cmd := exec.CommandContext(ctx, "wsl.exe", "--export", distroName, outputPath, "--format", "tar.gz")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("export failed: %v (output: %s)", err, string(output))
	}

	return nil
}

// BackupDistros backs up one or more distros
func BackupDistros(ctx context.Context, b Backuper, distros []string, backupDir string, opts BackupOptions) []BackupResult {
	results := make([]BackupResult, 0, len(distros))

	if len(distros) == 0 {
		return results
	}

	for _, distroName := range distros {
		result := BackupResult{
			Distro:  distroName,
			Success: false,
		}

		// Verify distro is registered
		registered, err := b.IsRegistered(ctx, distroName)
		if err != nil {
			result.Message = fmt.Sprintf("Error checking registration: %v", err)
			results = append(results, result)
			continue
		}
		if !registered {
			result.Message = "Distribution is not registered"
			results = append(results, result)
			continue
		}

		// Generate filename
		var filename string
		if opts.CustomName != "" {
			filename = opts.CustomName
			// Add extension if not present
			if filepath.Ext(filename) != ".gz" && filepath.Ext(filename) != ".tar" {
				filename += ".tar.gz"
			}
		} else {
			// Auto-generate name with timestamp
			timestamp := time.Now().Format("20060102-150405")
			filename = fmt.Sprintf("%s-%s.tar.gz", distroName, timestamp)
		}

		outputPath := filepath.Join(backupDir, filename)

		// Perform export
		err = b.Export(ctx, distroName, outputPath)
		if err != nil {
			result.Message = fmt.Sprintf("Backup failed: %v", err)
			results = append(results, result)
			continue
		}

		result.Success = true
		result.FilePath = outputPath
		result.Message = "Backup completed successfully"

		results = append(results, result)
	}

	return results
}
