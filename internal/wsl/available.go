package wsl

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// AvailableDistro represents a distro available for installation
type AvailableDistro struct {
	Name         string `json:"name"`
	FriendlyName string `json:"friendlyName"`
}

// decodeUTF16 converts UTF-16 bytes to UTF-8 string
func decodeUTF16(b []byte) string {
	if len(b)%2 != 0 {
		return string(b)
	}

	u16s := make([]uint16, 1)
	ret := &strings.Builder{}
	b8buf := make([]byte, 4)

	lb := len(b)
	for i := 0; i < lb; i += 2 {
		u16s[0] = uint16(b[i]) + (uint16(b[i+1]) << 8)
		r := utf16.Decode(u16s)
		n := utf8.EncodeRune(b8buf, r[0])
		ret.Write(b8buf[:n])
	}

	return ret.String()
}

// GetAvailableDistros retrieves the list of distros available to install
func GetAvailableDistros(ctx context.Context) ([]AvailableDistro, error) {
	cmd := exec.CommandContext(ctx, "wsl", "--list", "--online")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get available distros: %w", err)
	}

	// Decode UTF-16 to UTF-8
	text := decodeUTF16(output)
	
	var distros []AvailableDistro
	scanner := bufio.NewScanner(strings.NewReader(text))
	
	// Skip header lines
	headerCount := 0
	for scanner.Scan() {
		headerCount++
		if headerCount >= 3 {
			break
		}
	}

	// Parse distro lines
	for scanner.Scan() {
		lineText := scanner.Text()
		if strings.TrimSpace(lineText) == "" {
			continue
		}

		// Split by whitespace, taking first part as name and rest as friendly name
		fields := strings.Fields(lineText)
		if len(fields) >= 2 {
			name := fields[0]
			friendlyName := strings.Join(fields[1:], " ")
			
			// Skip the header row that contains column names
			if name == "NAME" && strings.Contains(friendlyName, "FRIENDLY") {
				continue
			}
			
			distros = append(distros, AvailableDistro{
				Name:         name,
				FriendlyName: friendlyName,
			})
		}
	}

	return distros, nil
}
