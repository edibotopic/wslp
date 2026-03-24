package wsl

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

const ubuntuRegistryKey = `Software\Canonical\Ubuntu`
const ubuntuInsightsConsent = "UbuntuInsightsConsent"

// GetUbuntuTelemetryStatus returns whether Ubuntu analytics consent is enabled.
func GetUbuntuTelemetryStatus() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, ubuntuRegistryKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	consent, _, err := key.GetIntegerValue(ubuntuInsightsConsent)
	if err != nil {
		return false
	}

	return consent == 1
}

// SetUbuntuTelemetryStatus sets the Ubuntu analytics consent registry value.
func SetUbuntuTelemetryStatus(enabled bool) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, ubuntuRegistryKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open Ubuntu registry key: %w", err)
	}
	defer key.Close()

	value := uint32(0)
	if enabled {
		value = 1
	}

	if err := key.SetDWordValue(ubuntuInsightsConsent, value); err != nil {
		return fmt.Errorf("failed to set %s: %w", ubuntuInsightsConsent, err)
	}

	return nil
}
