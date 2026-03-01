package wsl

import (
	"context"
	"fmt"
	"strings"

	gowsl "github.com/ubuntu/gowsl"
	"golang.org/x/sys/windows/registry"
)

// WSLSystemInfo contains system-wide WSL information
type WSLSystemInfo struct {
	DefaultWSLVersion int    `json:"defaultWslVersion"`
	NumDistros        int    `json:"numDistros"`
	TotalDiskUsage    string `json:"totalDiskUsage"` // Placeholder for now
	DefaultDistro     string `json:"defaultDistro"`
}

// DistroDetailInfo contains detailed information about a specific distro
type DistroDetailInfo struct {
	Name             string            `json:"name"`
	WSLVersion       int               `json:"wslVersion"`
	State            string            `json:"state"`
	IsDefault        bool              `json:"isDefault"`
	GUID             string            `json:"guid"`
	DefaultUID       uint32            `json:"defaultUid"`
	InteropEnabled   bool              `json:"interopEnabled"`
	DriveMounting    bool              `json:"driveMounting"`
	PathAppended     bool              `json:"pathAppended"`
	Flavor           string            `json:"flavor"`                     // e.g., "ubuntu", "debian"
	IsUbuntu         bool              `json:"isUbuntu"`
	TelemetryEnabled *bool             `json:"telemetryEnabled,omitempty"` // Only for Ubuntu
	EnvironmentVars  map[string]string `json:"environmentVars"`
}

// GetWSLSystemInfo retrieves system-wide WSL information
func GetWSLSystemInfo(ctx context.Context) (WSLSystemInfo, error) {
	info := WSLSystemInfo{}

	// Get all registered distros
	distros, err := gowsl.RegisteredDistros(ctx)
	if err != nil {
		return info, fmt.Errorf("failed to get distros: %w", err)
	}
	info.NumDistros = len(distros)

	// Get default distro
	defaultDistro, ok, err := gowsl.DefaultDistro(ctx)
	if err != nil {
		return info, fmt.Errorf("failed to get default distro: %w", err)
	}
	if ok {
		info.DefaultDistro = defaultDistro.Name()
	}

	// Get default WSL version from registry
	// HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Lxss\DefaultVersion
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Lxss`, registry.QUERY_VALUE)
	if err == nil {
		defer key.Close()
		version, _, err := key.GetIntegerValue("DefaultVersion")
		if err == nil {
			info.DefaultWSLVersion = int(version)
		} else {
			// Default to 2 if not set
			info.DefaultWSLVersion = 2
		}
	} else {
		info.DefaultWSLVersion = 2
	}

	// TODO: Calculate total disk usage
	info.TotalDiskUsage = "N/A"

	return info, nil
}

// GetDistroDetailInfo retrieves detailed information about a specific distro
func GetDistroDetailInfo(ctx context.Context, name string) (DistroDetailInfo, error) {
	info := DistroDetailInfo{
		Name: name,
	}

	distro := gowsl.NewDistro(ctx, name)

	// Check if registered
	registered, err := distro.IsRegistered()
	if err != nil {
		return info, fmt.Errorf("failed to check registration: %w", err)
	}
	if !registered {
		return info, fmt.Errorf("distro %s is not registered", name)
	}

	// Get state
	state, err := distro.State()
	if err == nil {
		info.State = state.String()
	}

	// Get GUID
	guid, err := distro.GUID()
	if err == nil {
		guidStr := guid.String()
		info.GUID = guidStr

		// Get Flavor from registry using GUID
		flavor, err := getDistroFlavor(guidStr)
		if err == nil {
			info.Flavor = flavor
			info.IsUbuntu = strings.EqualFold(flavor, "ubuntu")

			// If Ubuntu, check telemetry status
			if info.IsUbuntu {
				telemetry := getUbuntuTelemetryStatus()
				info.TelemetryEnabled = &telemetry
			}
		}
	}

	// Get configuration
	config, err := distro.GetConfiguration()
	if err == nil {
		info.WSLVersion = int(config.Version)
		info.DefaultUID = config.DefaultUID
		info.InteropEnabled = config.InteropEnabled
		info.DriveMounting = config.DriveMountingEnabled
		info.PathAppended = config.PathAppended
		info.EnvironmentVars = config.DefaultEnvironmentVariables
	}

	// Check if default - use case-insensitive comparison since gowsl may return lowercase
	// NOTE: Distros with matching names but different casings is atypical but possible.
	defaultDistro, ok, err := gowsl.DefaultDistro(ctx)
	if err == nil && ok {
		info.IsDefault = strings.EqualFold(defaultDistro.Name(), name)
	}

	return info, nil
}

// getDistroFlavor retrieves the Flavor value from the registry
func getDistroFlavor(guid string) (string, error) {
	// Ensure GUID has braces
	if guid[0] != '{' {
		guid = "{" + guid + "}"
	}

	keyPath := fmt.Sprintf(`Software\Microsoft\Windows\CurrentVersion\Lxss\%s`, guid)
	key, err := registry.OpenKey(registry.CURRENT_USER, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return "", fmt.Errorf("failed to open distro registry key: %w", err)
	}
	defer key.Close()

	flavor, _, err := key.GetStringValue("Flavor")
	if err != nil {
		return "", fmt.Errorf("failed to read Flavor value: %w", err)
	}

	return flavor, nil
}

// getUbuntuTelemetryStatus checks the Ubuntu telemetry consent registry key
func getUbuntuTelemetryStatus() bool {
	// HKEY_CURRENT_USER\Software\Canonical\Ubuntu\UbuntuInsightsConsent
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Canonical\Ubuntu`, registry.QUERY_VALUE)
	if err != nil {
		// Key doesn't exist, assume consent not given
		return false
	}
	defer key.Close()

	consent, _, err := key.GetIntegerValue("UbuntuInsightsConsent")
	if err != nil {
		return false
	}

	return consent == 1
}
