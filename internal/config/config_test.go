package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestDefaultBackupDir(t *testing.T) {
	t.Run("uses USERPROFILE when set", func(t *testing.T) {
		t.Setenv("USERPROFILE", `C:\Users\testuser`)
		got := DefaultBackupDir()
		want := filepath.Join(`C:\Users\testuser`, "WSLBackups")
		if got != want {
			t.Errorf("DefaultBackupDir() = %q, want %q", got, want)
		}
	})

	t.Run("falls back to home dir when USERPROFILE empty", func(t *testing.T) {
		t.Setenv("USERPROFILE", "")
		home, err := os.UserHomeDir()
		if err != nil {
			t.Skip("cannot determine home dir on this machine")
		}
		got := DefaultBackupDir()
		want := filepath.Join(home, "WSLBackups")
		if got != want {
			t.Errorf("DefaultBackupDir() = %q, want %q", got, want)
		}
	})
}

func TestSetDefaults(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	SetDefaults()

	if got := viper.GetInt("max_concurrent_installs"); got != 3 {
		t.Errorf("max_concurrent_installs default = %d, want 3", got)
	}

	if got := viper.GetString("backup_dir"); got == "" {
		t.Errorf("backup_dir default is empty, want a non-empty path")
	}
}

func TestGetMaxConcurrentInstalls(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	t.Run("returns configured value", func(t *testing.T) {
		viper.Set("max_concurrent_installs", 7)
		if got := GetMaxConcurrentInstalls(); got != 7 {
			t.Errorf("GetMaxConcurrentInstalls() = %d, want 7", got)
		}
	})

	t.Run("returns default after SetDefaults", func(t *testing.T) {
		viper.Reset()
		SetDefaults()
		if got := GetMaxConcurrentInstalls(); got != 3 {
			t.Errorf("GetMaxConcurrentInstalls() = %d, want 3", got)
		}
	})
}

func TestGetBackupDir(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	viper.Set("backup_dir", `C:\custom\backup\dir`)
	if got := GetBackupDir(); got != `C:\custom\backup\dir` {
		t.Errorf("GetBackupDir() = %q, want %q", got, `C:\custom\backup\dir`)
	}
}

func TestEnsureBackupDir(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	dir := filepath.Join(t.TempDir(), "nested", "backup", "dir")
	viper.Set("backup_dir", dir)

	if err := EnsureBackupDir(); err != nil {
		t.Fatalf("EnsureBackupDir() error = %v", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("expected backup dir to be created, stat error: %v", err)
	}
	if !info.IsDir() {
		t.Errorf("expected %q to be a directory", dir)
	}
}

func TestInit(t *testing.T) {
	viper.Reset()
	defer viper.Reset()

	// Init should not panic and should populate defaults even if no
	// config file is present.
	Init()

	if got := viper.GetInt("max_concurrent_installs"); got != 3 {
		t.Errorf("after Init(), max_concurrent_installs = %d, want 3", got)
	}
	if got := viper.GetString("backup_dir"); got == "" {
		t.Errorf("after Init(), backup_dir is empty, want a non-empty path")
	}
}
