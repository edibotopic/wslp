## wslp backup

Backup one or more WSL distributions

### Synopsis

Backup one or more WSL distributions to tar.gz files.

By default, backups are saved to %USERPROFILE%\WSLBackups with an auto-generated
name including the distro name and timestamp (e.g., Ubuntu-20240301-143022.tar.gz).

You can specify a custom name for single distro backups using the --name flag.
The backup directory can be customized via the --backup-dir flag or by setting
backup_dir in ~/.wslp.yaml, or via the WSLP_BACKUP_DIR environment variable.

```
wslp backup <distro> [distro...] [flags]
```

### Options

```
  -d, --backup-dir string   Directory to save backups (overrides config)
  -h, --help                help for backup
  -n, --name string         Custom name for the backup file (only for single distro)
```

### SEE ALSO

* [wslp](wslp.md)	 - A tool for managing WSL instances.

