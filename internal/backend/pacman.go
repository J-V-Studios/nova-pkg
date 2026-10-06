package backend

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Pacman is the backend for the system package manager.
type Pacman struct{}

func (Pacman) Name() string { return "pacman" }

func (Pacman) Available() bool {
	_, err := exec.LookPath("pacman")
	return err == nil
}

func (Pacman) Search(query string) ([]Package, error) {
	out, err := exec.Command("pacman", "-Ss", query).Output()
	if err != nil {
		// pacman exits 1 when nothing matches; that is not an error for us.
		if len(out) == 0 {
			return nil, nil
		}
		return nil, fmt.Errorf("pacman search: %w", err)
	}
	return parseSearch(string(out)), nil
}

func (Pacman) Install(pkg string) error {
	cmd := exec.Command("sudo", "pacman", "-S", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pacman install %s: %w", pkg, err)
	}
	return nil
}

// parseSearch parses `pacman -Ss` output. Pure function so it can be
// unit tested without running pacman.
//
// Format: a header line "repo/name version [flags]" followed by an
// indented description line.
func parseSearch(output string) []Package {
	var pkgs []Package
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if line == "" || line[0] == ' ' || line[0] == '\t' {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		name := fields[0]
		if idx := strings.IndexByte(name, '/'); idx >= 0 {
			name = name[idx+1:] // strip "repo/" prefix
		}
		desc := ""
		if i+1 < len(lines) && (strings.HasPrefix(lines[i+1], " ") || strings.HasPrefix(lines[i+1], "\t")) {
			desc = strings.TrimSpace(lines[i+1])
		}
		pkgs = append(pkgs, Package{
			Name:        name,
			Version:     fields[1],
			Description: desc,
			Source:      "pacman",
		})
	}
	return pkgs
}
