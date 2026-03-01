# Building

The app is experimental and was made for prototyping purposes.

It is currently not available as a packaged release.

This guide covers building both the CLI and GUI components.

## Prerequisites

- **Go** 1.20 or later
- **Flutter** 3.0 or later
- **Windows** 11

## Building the CLI

The CLI is a Go application using Cobra for command-line interface.

You first need to clone the repo:

```bash
git clone https://github.com/edibotopic/wslp
```

### Standard build

```bash
cd wslp
go build -o wslp.exe
```

This creates an executable `wslp.exe` in the current directory.

### (optional) Install globally

```bash
go install
```

This installs `wslp` to your `$GOPATH/bin` directory.

## Building the GUI

The GUI is a Flutter application that communicates with the Go backend
via HTTP.

### Navigate to GUI directory

```bash
cd gui
```

### Build for Windows

```bash
flutter build windows
```
The compiled application will be in `gui/build/windows/runner/Release/`.

## Using the GUI

The GUI requires the backend server to be running.

Start it with:

```bash
wslp.exe serve
```

Then run the GUI.

## Development

### Run CLI in development

```bash
go run main.go [command]
```

Example:

```bash
go run main.go list
```

### Run GUI in development

```bash
cd gui
flutter run -d windows
```
## Testing

### Test Go code

```bash
go test ./...
```

### Test with coverage

```bash
go test -cover ./...
```

### Test Specific Package

```bash
go test ./internal/wsl
go test ./cmd
```

## Cross-compilation

Since WSL Plus is Windows-specific, cross-compilation to other platforms
is not supported.
