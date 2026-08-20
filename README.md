# WSL Plus

Prototype management tool for WSL.

[docs-image]: https://readthedocs.org/projects/wslp/badge/?version=latest
[docs-url-stable]: https://wslp.readthedocs.io/

## Usage

WSL Plus provides two interfaces for managing WSL distributions: a command-line interface (CLI) and a graphical user interface (GUI).

### Prerequisites for building

- [Go](https://go.dev/dl/) — required to build the CLI
- [Flutter](https://docs.flutter.dev/install) — required to build the GUI

With `go` and `flutter` installed and on your PATH, you can then build
the CLI and the tool using the batch scripts in the repo.

### Quickstart

After cloning the repo, build the CLI and the GUI with `build.bat`.

Then, still within root of the repo, the following steps will enable usage of
the CLI and GUI.

- **For CLI**: Use `wslp.exe` directly

  ```bash
  .\wslp.exe list
  .\wslp.exe backup Ubuntu
  ```

- **For GUI**: Run `rungui.bat`

This last command automatically starts the backend server and launches the GUI.

#### Installing the CLI tool to your path

From within the repo's root, run:

```bash
go install .
```

Now you can call `wslp` directly from anywhere.

### CLI Usage

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
wslp install Ubuntu Debian archlinux
```

There is also a server that is used as the backend for the GUI.

```bash
wslp serve
```

This starts the HTTP API server on port 8080 (default). This is required
for the GUI to function.

### GUI Usage

The GUI provides a visual interface for managing WSL distributions.

The main screen displays all registered WSL distributions as cards. The
default distribution is highlighted with a badge and appears first in
the list.

The GUI offers both bulk commands, which are accessible from the side
navigation, and per-distro commands, which are accessible in the context
menu for each distro.
