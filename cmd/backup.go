package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"wslp/internal/config"
	"wslp/internal/wsl"
)

func init() {
	RootCmd.AddCommand(newBackupCmd())
}

func newBackupCmd() *cobra.Command {
	var customName string
	var backupDir string

	cmd := &cobra.Command{
		Use:   "backup <distro> [distro...]",
		Short: "Backup one or more WSL distributions",
		Long: `Backup one or more WSL distributions to tar.gz files.

By default, backups are saved to %USERPROFILE%\WSLBackups with an auto-generated
name including the distro name and timestamp (e.g., Ubuntu-20240301-143022.tar.gz).

You can specify a custom name for single distro backups using the --name flag.
The backup directory can be customized via the --backup-dir flag or by setting
backup_dir in ~/.wslp.yaml, or via the WSLP_BACKUP_DIR environment variable.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBackup(cmd.OutOrStdout(), args, customName, backupDir)
		},
	}

	cmd.Flags().StringVarP(&customName, "name", "n", "", "Custom name for the backup file (only for single distro)")
	cmd.Flags().StringVarP(&backupDir, "backup-dir", "d", "", "Directory to save backups (overrides config)")

	return cmd
}

func runBackup(w io.Writer, distros []string, customName, backupDir string) error {
	ctx := context.Background()

	// If customName is provided but multiple distros, return error
	if customName != "" && len(distros) > 1 {
		return fmt.Errorf("custom name can only be used when backing up a single distribution")
	}

	// Determine backup directory
	if backupDir == "" {
		backupDir = config.GetBackupDir()
	}

	// Ensure backup directory exists
	if err := config.EnsureBackupDir(); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	opts := wsl.BackupOptions{
		CustomName: customName,
	}

	backuper := wsl.RealBackuper{}
	results := wsl.BackupDistros(ctx, backuper, distros, backupDir, opts)

	// Print results
	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
			fmt.Fprintf(w, "✓ %s: %s\n", result.Distro, result.Message)
			fmt.Fprintf(w, "  Saved to: %s\n", result.FilePath)
		} else {
			fmt.Fprintf(w, "✗ %s: %s\n", result.Distro, result.Message)
		}
	}

	if successCount > 0 {
		fmt.Fprintf(w, "\nSuccessfully backed up %d/%d distribution(s) to %s\n", successCount, len(results), backupDir)
	}

	if successCount < len(results) {
		return fmt.Errorf("some backups failed")
	}

	return nil
}
