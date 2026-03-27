## Why wslp

Examples of why wslp is useful.

### Streamlined actions

#### Renaming

Renaming a distro requires editing the Windows Registry directly when using `wsl.exe`.

::::{grid} 2

:::{grid-item-card} Without wslp
1. Open Registry Editor (`regedit`)
2. Navigate to `HKEY_CURRENT_USER\SOFTWARE\Microsoft\Windows\CurrentVersion\Lxss`
3. Find the subkey for your distro by checking each `DistributionName` value
4. Double-click `DistributionName` and enter the new name
5. Run `wsl --shutdown` to apply the change
:::

:::{grid-item-card} With wslp
1. Run `wslp rename <old-name> <new-name>`
2. Run `wsl --shutdown` to apply the change
:::

::::

#### Copying

Copying normally requires at least two steps: exporting and importing:

::::{grid} 2

:::{grid-item-card} Without wslp
1. mkdir backup
2. `wsl --export Ubuntu-24.04 .\backup\Ubuntu-24.04.tar.gz`
3. `wsl --import Ubuntu-24.04-copy .\backup\Ubuntu2404-copy\ .\backup\Ubuntu-24.04.tar.gz`
:::

:::{grid-item-card} With wslp
1. Run `wslp copy Ubuntu-24.04 Ubuntu-24.04-copy`
:::

::::


### Bulk actions

Installing multiple distros with `wsl.exe` requires a separate command per distro.

::::{grid} 2

:::{grid-item-card} Without wslp
1. Run `wsl --install Ubuntu`
2. Wait for it to complete
3. Run `wsl --install Debian`
4. Wait for it to complete
5. Repeat for each additional distro
:::

:::{grid-item-card} With wslp
1. Run `wslp install <distro1> <distro2> ...`.
:::

::::

```{note}
With [concurrent installations](./concurrent-installs), you can also reduce the time for installation.
```

Backing up multiple distros with `wsl.exe` also requires a separate command per distro.

::::{grid} 2

:::{grid-item-card} Without wslp
1. Run `wsl --export Ubuntu Ubuntu.tar`
2. Wait for it to complete
3. Run `wsl --export Debian Debian.tar`
4. Wait for it to complete
5. Repeat for each additional distro
:::

:::{grid-item-card} With wslp
1. Run `wslp backup Ubuntu Debian ...`
:::

::::

### Context and feedback

Getting an overview of all distros and inspecting their properties with
`wsl.exe` requires running multiple commands, reviewing separate outputs. For
specific inspections, like Ubuntu telemetry, the registry is required.

::::{grid} 2

:::{grid-item-card} Without wslp
1. Run `wsl --list --verbose` to see distro names and running state
2. Run `wsl -d Ubuntu` to enter a distro's shell
3. Run distro-specific commands to inspect properties (e.g., check Ubuntu telemetry settings)
4. Exit the shell and repeat from step 2 for each other distro
:::

:::{grid-item-card} With wslp
1. Start the GUI
2. View all distros and their state at a glance
3. Select a distro's context menu to inspect its properties
:::

::::

### Feature comparison

| Feature | `wsl.exe` | `wslp` |
|---------|:---------:|:------:|
| Install a distro | ✓ | ✓ |
| Install multiple distros at once | ✗ | ✓ |
| List distros | ✓ | ✓ |
| Launch a distro shell | ✓ | ✓ |
| Terminate a distro | ✓ | ✓ |
| Terminate multiple distros at once | ✗ | ✓ |
| Unregister a distro | ✓ | ✓ |
| Unregister multiple distros at once | ✗ | ✓ |
| Backup a distro | ✓ | ✓ |
| Backup multiple distros at once | ✗ | ✓ |
| Auto backup filename and location | ✗ | ✓ |
| Rename a distro | ✗ | ✓ |
| GUI overview of distros and their state | ✗ | ✓ |
| Inspect per-distro properties | ✗ | ✓ |
