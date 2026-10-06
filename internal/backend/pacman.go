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
	// pacman -Ss treats the query as a regex; only flag injection
	// (leading dash) must be blocked, not regex syntax itself.
	if strings.HasPrefix(query, "-") {
		return nil, fmt.Errorf("invalid query %q", query)
	}
	out, err := exec.Command("pacman", "-Ss", "--", query).Output()
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
	if err := ValidateName(pkg); err != nil {
		return err
	}
	cmd := exec.Command("sudo", "pacman", "-S", "--", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pacman install %s: %w", pkg, err)
	}
	return nil
}

func (Pacman) Remove(pkg string) error {
	if err := ValidateName(pkg); err != nil {
		return err
	}
	cmd := exec.Command("sudo", "pacman", "-R", "--", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pacman remove %s: %w", pkg, err)
	}
	return nil
}

func (Pacman) Update() error {
	cmd := exec.Command("sudo", "pacman", "-Syu")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pacman update: %w", err)
	}
	return nil
}

func (Pacman) List() ([]Package, error) {
	out, err := exec.Command("pacman", "-Q").Output()
	if err != nil {
		return nil, fmt.Errorf("pacman list: %w", err)
	}
	return parsePacmanQ(string(out), "pacman"), nil
}

// parsePacmanQ parses `pacman -Q` / `pacman -Qm` output: one
// "name version" line per installed package. Pure function.
func parsePacmanQ(output, source string) []Package {
	var pkgs []Package
	for _, line := range strings.Split(output, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		pkgs = append(pkgs, Package{
			Name:    fields[0],
			Version: fields[1],
			Source:  source,
		})
	}
	return pkgs
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
