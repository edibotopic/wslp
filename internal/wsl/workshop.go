package wsl

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// WorkshopInfo represents a single Canonical Workshop (canonical/workshop)
// development environment, as reported by `workshop list --global` inside
// a WSL distro.
type WorkshopInfo struct {
	Project string `json:"project"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

// WorkshopRunner executes `workshop list` inside a distro and returns its
// raw stdout. Abstracted so the parsing logic can be unit tested without
// shelling out.
type WorkshopRunner interface {
	ListWorkshops(ctx context.Context, distro string) ([]byte, error)
}

// RealWorkshopRunner shells out to `wsl -d <distro> -- workshop list ...`.
type RealWorkshopRunner struct{}

func (RealWorkshopRunner) ListWorkshops(ctx context.Context, distro string) ([]byte, error) {
	// Workshop is typically installed as a snap (/snap/bin/workshop), whose
	// directory is only added to PATH by login-shell profile scripts (e.g.
	// /etc/profile.d/apps-bin-path.sh). `wsl -d <distro> -- workshop ...`
	// runs a non-login shell, so PATH lacks /snap/bin and the command fails
	// with "command not found". Run it through `bash -lc` instead so login
	// profile scripts run and PATH is set up correctly.
	cmd := exec.CommandContext(ctx, "wsl", "-d", distro, "--", "bash", "-lc", "workshop list --global --no-headers")
	return cmd.Output()
}

// GetWorkshops lists the Workshop environments running inside distro.
//
// Workshop (https://github.com/canonical/workshop) is an optional,
// third-party dev-environment tool that may not be installed in a given
// distro. Any failure (missing binary, daemon not running, distro
// unreachable) is treated as "no workshops" rather than a hard error,
// since its absence is expected and not a wslp failure.
func GetWorkshops(ctx context.Context, distro string, r WorkshopRunner) []WorkshopInfo {
	output, err := r.ListWorkshops(ctx, distro)
	if err != nil {
		return []WorkshopInfo{}
	}
	return parseWorkshopList(string(output))
}

// parseWorkshopList parses the plain-text table produced by
// `workshop list --global --no-headers`, with columns:
// Project, Workshop, Status, Notes (notes may be absent or contain spaces).
func parseWorkshopList(output string) []WorkshopInfo {
	workshops := []WorkshopInfo{}

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		workshops = append(workshops, WorkshopInfo{
			Project: fields[0],
			Name:    fields[1],
			Status:  fields[2],
		})
	}

	return workshops
}

// WorkshopController runs `workshop start|stop` inside a distro for a
// specific project/workshop and returns its combined output (used as the
// error message on failure). Abstracted so callers can be unit tested
// without shelling out.
type WorkshopController interface {
	RunAction(ctx context.Context, distro, project, name, action string) ([]byte, error)
}

// RealWorkshopController shells out to
// `wsl -d <distro> -- bash -lc "workshop <action> <name> -p <project>"`.
type RealWorkshopController struct{}

func (RealWorkshopController) RunAction(ctx context.Context, distro, project, name, action string) ([]byte, error) {
	command := fmt.Sprintf("workshop %s %s -p %s", action, shellQuote(name), shellProjectArg(project))
	cmd := exec.CommandContext(ctx, "wsl", "-d", distro, "--", "bash", "-lc", command)
	return cmd.CombinedOutput()
}

// StartWorkshop activates the named workshop (must have been launched
// previously). See `workshop start --help`.
func StartWorkshop(ctx context.Context, distro, project, name string, c WorkshopController) error {
	return runWorkshopControlAction(ctx, distro, project, name, "start", c)
}

// StopWorkshop deactivates the named workshop. See `workshop stop --help`.
func StopWorkshop(ctx context.Context, distro, project, name string, c WorkshopController) error {
	return runWorkshopControlAction(ctx, distro, project, name, "stop", c)
}

func runWorkshopControlAction(ctx context.Context, distro, project, name, action string, c WorkshopController) error {
	output, err := c.RunAction(ctx, distro, project, name, action)
	if err != nil {
		msg := strings.TrimSpace(string(output))
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("workshop %s failed: %s", action, msg)
	}
	return nil
}

// LaunchWorkshopShell opens an interactive `workshop shell` session for the
// named workshop in a new terminal window (non-blocking), mirroring
// LaunchInTerminal's approach for plain WSL distros. The workshop must be
// "Ready" or "Waiting" for the shell command to succeed.
func LaunchWorkshopShell(ctx context.Context, distro, project, name string) error {
	command := fmt.Sprintf("workshop shell %s -p %s", shellQuote(name), shellProjectArg(project))

	if isWindowsTerminalAvailable() {
		cmd := exec.CommandContext(ctx, "wt.exe", "wsl.exe", "-d", distro, "--", "bash", "-lc", command)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to launch Windows Terminal: %w", err)
		}
		return nil
	}

	cmd := exec.CommandContext(ctx, "wsl.exe", "-d", distro, "--", "bash", "-lc", command)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch wsl.exe: %w", err)
	}
	return nil
}

// shellQuote wraps s in single quotes for safe use inside a `bash -lc`
// command string, escaping any embedded single quotes.
func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// shellProjectArg builds a safely-quoted bash word for a project path that
// may use the `~` (home directory) shorthand, as reported by
// `workshop list --global`. Single-quoting the whole path (via shellQuote)
// would suppress bash's tilde expansion, causing `workshop` to look for a
// literal "~" directory and fail. Instead, a leading "~" is rewritten to an
// unquoted "$HOME" that bash expands, concatenated with the safely-quoted
// remainder (adjacent quoted/unquoted words with no space form one bash
// word, so this still yields a single argument to `workshop`).
func shellProjectArg(project string) string {
	if project == "~" {
		return `"$HOME"`
	}
	if rest, ok := strings.CutPrefix(project, "~/"); ok {
		return `"$HOME"` + shellQuote("/"+rest)
	}
	return shellQuote(project)
}
