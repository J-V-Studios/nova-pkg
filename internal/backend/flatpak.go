package backend

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Flatpak is the backend for Flatpak remotes such as Flathub.
type Flatpak struct{}

func (Flatpak) Name() string { return "flatpak" }

func (Flatpak) Available() bool {
	_, err := exec.LookPath("flatpak")
	return err == nil
}

func (Flatpak) Search(query string) ([]Package, error) {
	out, err := exec.Command("flatpak", "search",
		"--columns=name,description,application,version", query).Output()
	if err != nil {
		if len(out) == 0 {
			return nil, nil // no matches
		}
		return nil, fmt.Errorf("flatpak search: %w", err)
	}
	return parseFlatpakSearch(string(out)), nil
}

func (Flatpak) Install(pkg string) error {
	// -y answers yes to prompts; stdin stays attached so remote
	// selection and polkit auth still work when -y cannot cover them.
	cmd := exec.Command("flatpak", "install", "-y", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("flatpak install %s: %w", pkg, err)
	}
	return nil
}

// parseFlatpakSearch parses `flatpak search --columns=name,description,application,version`
// output: tab-separated rows. flatpak omits the header when stdout is
// not a terminal, so rows are identified by field count, not position.
// Pure function so it can be unit tested without running flatpak.
func parseFlatpakSearch(output string) []Package {
	var pkgs []Package
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}
		cols := strings.Split(line, "\t")
		if len(cols) < 4 {
			continue // malformed line
		}
		// On a terminal flatpak prints a header row; skip it. Versions
		// always contain a digit or "r<digits>", headers never do.
		if !strings.ContainsAny(cols[3], "0123456789") {
			continue
		}
		pkgs = append(pkgs, Package{
			Name:        cols[2], // application ID, used by flatpak install
			Version:     cols[3],
			Description: cols[0] + " - " + cols[1],
			Source:      "flatpak",
		})
	}
	return pkgs
}
