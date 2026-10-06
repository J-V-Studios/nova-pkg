// Package backend defines the common interface every package source
// (pacman, AUR, Flatpak, ...) implements.
package backend

import (
	"fmt"
	"regexp"
)

// Package is a single result from any package source.
type Package struct {
	Name        string
	Version     string
	Description string
	Source      string
}

// Backend wraps one package source behind a common interface.
type Backend interface {
	Name() string
	Available() bool // exec.LookPath check for the tool
	Search(query string) ([]Package, error)
	Install(pkg string) error
	Remove(pkg string) error
	Update() error
	List() ([]Package, error) // installed packages
}

// validNameRe rejects flag injection (leading dash), path traversal
// (/, ..), and shell metacharacters in package names. Covers pacman,
// AUR (repo names are [a-z0-9@._+-]), and flatpak app IDs
// (reverse-DNS dotted names).
var validNameRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9@._+-]*$`)

// ValidateName rejects package names that could be interpreted as
// flags or paths by the tools backends shell out to.
func ValidateName(name string) error {
	if !validNameRe.MatchString(name) {
		return fmt.Errorf("invalid package name %q", name)
	}
	return nil
}
