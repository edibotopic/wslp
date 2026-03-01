package wsl

import (
	"context"
	"fmt"

	gowsl "github.com/ubuntu/gowsl"
	"golang.org/x/sys/windows/registry"
)

// RenameResult contains the result of renaming a distro
type RenameResult struct {
	OldName string `json:"oldName"`
	NewName string `json:"newName"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Renamer interface for renaming distros
type Renamer interface {
	IsRegistered(ctx context.Context, name string) (bool, error)
	GetDistroGUID(ctx context.Context, name string) (string, error)
	RenameInRegistry(guid, newName string) error
}

// RealRenamer implements Renamer using Windows Registry
type RealRenamer struct{}

// IsRegistered checks if a distro is registered
func (r RealRenamer) IsRegistered(ctx context.Context, name string) (bool, error) {
	distros, err := gowsl.RegisteredDistros(ctx)
	if err != nil {
		return false, err
	}
	// NOTE: Exact match comparison. Distros with matching names but different
	// casings (e.g., "Ubuntu" and "ubuntu") is atypical but possible.
	for _, d := range distros {
		if d.Name() == name {
			return true, nil
		}
	}
	return false, nil
}

// GetDistroGUID gets the GUID for a distro using gowsl
func (r RealRenamer) GetDistroGUID(ctx context.Context, name string) (string, error) {
	distro := gowsl.NewDistro(ctx, name)

	guid, err := distro.GUID()
	if err != nil {
		return "", fmt.Errorf("failed to get GUID for distro %s: %w", name, err)
	}

	// Registry keys use GUIDs with braces: {GUID}
	// Ensure the GUID string has braces
	guidStr := guid.String()
	if guidStr[0] != '{' {
		guidStr = "{" + guidStr + "}"
	}

	return guidStr, nil
}

// RenameInRegistry updates the DistributionName value in the registry
func (r RealRenamer) RenameInRegistry(guid, newName string) error {
	// Open the specific distro's registry key
	keyPath := fmt.Sprintf(`Software\Microsoft\Windows\CurrentVersion\Lxss\%s`, guid)
	key, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open distro registry key: %w", err)
	}
	defer key.Close()

	// Update the DistributionName value
	err = key.SetStringValue("DistributionName", newName)
	if err != nil {
		return fmt.Errorf("failed to set new name: %w", err)
	}

	return nil
}

// RenameDistro renames a WSL distro by modifying the Windows Registry
func RenameDistro(ctx context.Context, r Renamer, oldName, newName string) RenameResult {
	result := RenameResult{
		OldName: oldName,
		NewName: newName,
		Success: false,
	}

	// Validate new name is not empty
	if newName == "" {
		result.Message = "New name cannot be empty"
		return result
	}

	// Check if old distro exists
	registered, err := r.IsRegistered(ctx, oldName)
	if err != nil {
		result.Message = fmt.Sprintf("Error checking registration: %v", err)
		return result
	}
	if !registered {
		result.Message = fmt.Sprintf("Distro %s is not registered", oldName)
		return result
	}

	// Check if new name already exists
	exists, err := r.IsRegistered(ctx, newName)
	if err != nil {
		result.Message = fmt.Sprintf("Error checking new name: %v", err)
		return result
	}
	if exists {
		result.Message = fmt.Sprintf("Distro %s already exists", newName)
		return result
	}

	// Get the GUID for the distro
	guid, err := r.GetDistroGUID(ctx, oldName)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to get distro GUID: %v", err)
		return result
	}

	// Rename in registry
	err = r.RenameInRegistry(guid, newName)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to rename: %v", err)
		return result
	}

	result.Success = true
	result.Message = fmt.Sprintf("Successfully renamed %s to %s. Please restart WSL for changes to take effect.", oldName, newName)

	return result
}
