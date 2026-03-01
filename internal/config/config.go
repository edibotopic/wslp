package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Init initializes the configuration with sensible defaults
func Init() {
	// Set config file details
	viper.SetConfigName(".wslp")
	viper.SetConfigType("yaml")

	// Look for config in user's home directory
	home, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(home)
	}

	// Also look in current directory
	viper.AddConfigPath(".")

	// Set defaults
	SetDefaults()

	// Read config file if it exists (ignore error if not found)
	viper.ReadInConfig()

	// Enable environment variable support with WSLP_ prefix
	viper.SetEnvPrefix("WSLP")
	viper.AutomaticEnv()
}

// SetDefaults sets default configuration values
func SetDefaults() {
	viper.SetDefault("backup_dir", DefaultBackupDir())
}

// DefaultBackupDir returns the default backup directory path
// Uses %USERPROFILE%\WSLBackups on Windows
func DefaultBackupDir() string {
	userProfile := os.Getenv("USERPROFILE")
	if userProfile == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			userProfile = home
		}
	}
	return filepath.Join(userProfile, "WSLBackups")
}

// GetBackupDir returns the configured backup directory
func GetBackupDir() string {
	return viper.GetString("backup_dir")
}

// EnsureBackupDir creates the backup directory if it doesn't exist
func EnsureBackupDir() error {
	backupDir := GetBackupDir()
	return os.MkdirAll(backupDir, 0755)
}
