## wslp rename

Rename a WSL distribution

### Synopsis

Rename a WSL distribution by modifying the Windows Registry.

This command updates the distribution name in the registry. After renaming,
you should restart WSL for the changes to take effect:

    wsl --shutdown

The rename operation:
- Validates the old distro exists
- Checks the new name doesn't conflict with existing distros
- Updates the registry entry directly (fast, no export/import needed)

```
wslp rename <old-name> <new-name> [flags]
```

### Options

```
  -h, --help   help for rename
```

### SEE ALSO

* [wslp](wslp.md)	 - A tool for managing WSL instances.

