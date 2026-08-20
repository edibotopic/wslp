package wsl

import (
	"testing"
)

func TestGetDistroFlavorWithBogusGUID(t *testing.T) {
	t.Run("returns error for non-existent GUID without braces", func(t *testing.T) {
		// Use a GUID that almost certainly won't exist in the registry
		guid := "00000000-0000-0000-0000-000000000000"

		flavor, err := getDistroFlavor(guid)

		if err == nil {
			t.Errorf("expected error for non-existent GUID, got flavor: %s", flavor)
		}
		if flavor != "" {
			t.Errorf("expected empty flavor on error, got: %s", flavor)
		}
	})

	t.Run("returns error for non-existent GUID with braces", func(t *testing.T) {
		// Use a GUID with braces that almost certainly won't exist
		guid := "{00000000-0000-0000-0000-000000000000}"

		flavor, err := getDistroFlavor(guid)

		if err == nil {
			t.Errorf("expected error for non-existent GUID, got flavor: %s", flavor)
		}
		if flavor != "" {
			t.Errorf("expected empty flavor on error, got: %s", flavor)
		}
	})

	t.Run("handles GUID normalization consistently for both forms", func(t *testing.T) {
		// Both forms of the same non-existent GUID should produce errors
		guidNoBraces := "ffffffff-ffff-ffff-ffff-ffffffffffff"
		guidWithBraces := "{ffffffff-ffff-ffff-ffff-ffffffffffff}"

		err1 := getDistroFlavorError(guidNoBraces)
		err2 := getDistroFlavorError(guidWithBraces)

		// Both should fail (as the GUID doesn't exist)
		if err1 == nil {
			t.Errorf("expected error for GUID without braces")
		}
		if err2 == nil {
			t.Errorf("expected error for GUID with braces")
		}
	})
}

func TestGetDistroRegistryVersionAndUIDWithBogusGUID(t *testing.T) {
	t.Run("returns error for non-existent GUID without braces", func(t *testing.T) {
		guid := "00000000-0000-0000-0000-000000000000"

		version, uid, err := getDistroRegistryVersionAndUID(guid)

		if err == nil {
			t.Errorf("expected error for non-existent GUID, got version: %d, uid: %d", version, uid)
		}
		if version != 0 || uid != 0 {
			t.Errorf("expected zero values on error, got version: %d, uid: %d", version, uid)
		}
	})

	t.Run("returns error for non-existent GUID with braces", func(t *testing.T) {
		guid := "{00000000-0000-0000-0000-000000000000}"

		version, uid, err := getDistroRegistryVersionAndUID(guid)

		if err == nil {
			t.Errorf("expected error for non-existent GUID, got version: %d, uid: %d", version, uid)
		}
		if version != 0 || uid != 0 {
			t.Errorf("expected zero values on error, got version: %d, uid: %d", version, uid)
		}
	})

	t.Run("handles GUID normalization consistently for both forms", func(t *testing.T) {
		// Both forms of the same non-existent GUID should produce errors
		guidNoBraces := "ffffffff-ffff-ffff-ffff-ffffffffffff"
		guidWithBraces := "{ffffffff-ffff-ffff-ffff-ffffffffffff}"

		err1 := getDistroRegistryVersionAndUIDError(guidNoBraces)
		err2 := getDistroRegistryVersionAndUIDError(guidWithBraces)

		// Both should fail (as the GUID doesn't exist)
		if err1 == nil {
			t.Errorf("expected error for GUID without braces")
		}
		if err2 == nil {
			t.Errorf("expected error for GUID with braces")
		}
	})
}

// Helper functions to check errors without requiring error values
func getDistroFlavorError(guid string) error {
	_, err := getDistroFlavor(guid)
	return err
}

func getDistroRegistryVersionAndUIDError(guid string) error {
	_, _, err := getDistroRegistryVersionAndUID(guid)
	return err
}
