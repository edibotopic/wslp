# Usage

WSL Plus provides two interfaces for managing WSL distributions: a command-line interface (CLI) and a graphical user interface (GUI).

## CLI Usage

The CLI provides direct command-line access to WSL management functions.

For example, to list distributions:

```bash
wslp list
```
Example output:

```
Found 2 registered distributions:
  - Ubuntu
  - Debian
```

Some commands can allow bulk operations, for example:

```bash
wslp install <distro-name> [distro-name...]
```

Examples:

```bash
# Install a single distribution
wslp install Ubuntu

# Install multiple distributions
wslp install Ubuntu Debian archos
```

There is also a server that is used as the backend for the GUI.

```bash
wslp serve
```

This starts the HTTP API server on port 8080 (default). This is required
for the GUI to function.

## GUI Usage

The GUI provides a visual interface for managing WSL distributions.

The main screen displays all registered WSL distributions as cards. The
default distribution is highlighted with a badge and appears first in
the list.

Click the **Refresh List** button in the sidebar to reload the current list of distributions.
Note that the list auto-refreshes every 5 seconds.

The GUI offers both bulk commands, which are accessible from the side
navigation, and per-distro commands, which are accessible in the context
menu for each distro.

The activity log at the bottom of the screen shows:
- Successful operations (marked with ✓ in green)
- Errors (marked with ✗ in red)
- Installation progress

Click **Clear Log** to reset the activity log.

