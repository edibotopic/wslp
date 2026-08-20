package wsl

import (
	"testing"
)

func TestDecodeUTF16(t *testing.T) {
	t.Run("handles empty input", func(t *testing.T) {
		result := decodeUTF16([]byte{})
		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})

	t.Run("handles simple ASCII UTF-16LE input", func(t *testing.T) {
		// "Hello" in UTF-16LE: H(0x48 0x00) e(0x65 0x00) l(0x6C 0x00) l(0x6C 0x00) o(0x6F 0x00)
		input := []byte{0x48, 0x00, 0x65, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0x6F, 0x00}
		result := decodeUTF16(input)
		if result != "Hello" {
			t.Errorf("expected 'Hello', got %q", result)
		}
	})

	t.Run("handles odd-length input by returning as-is", func(t *testing.T) {
		// Odd length should trigger fallback that returns string(b) unchanged
		input := []byte{0x48, 0x00, 0x65}
		result := decodeUTF16(input)
		if result != string(input) {
			t.Errorf("expected input returned as-is, got %q", result)
		}
	})

	t.Run("handles UTF-16LE with non-ASCII character (©)", func(t *testing.T) {
		// The copyright symbol © is U+00A9
		// In UTF-16LE: 0xA9 0x00
		// "©test" in UTF-16LE: © t e s t
		input := []byte{
			0xA9, 0x00, // © (U+00A9)
			0x74, 0x00, // t
			0x65, 0x00, // e
			0x73, 0x00, // s
			0x74, 0x00, // t
		}
		result := decodeUTF16(input)
		if result != "©test" {
			t.Errorf("expected '©test', got %q", result)
		}
	})

	t.Run("handles UTF-16LE with emoji-like multi-byte sequence", func(t *testing.T) {
		// Euro sign € is U+20AC
		// In UTF-16LE: 0xAC 0x20
		// "€100" in UTF-16LE
		input := []byte{
			0xAC, 0x20, // € (U+20AC)
			0x31, 0x00, // 1
			0x30, 0x00, // 0
			0x30, 0x00, // 0
		}
		result := decodeUTF16(input)
		if result != "€100" {
			t.Errorf("expected '€100', got %q", result)
		}
	})

	t.Run("handles single UTF-16 character", func(t *testing.T) {
		// Single character 'A' in UTF-16LE: 0x41 0x00
		input := []byte{0x41, 0x00}
		result := decodeUTF16(input)
		if result != "A" {
			t.Errorf("expected 'A', got %q", result)
		}
	})
}
