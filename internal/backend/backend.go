package backend

// Package backend defines the common interface every package source
// (pacman, AUR, Flatpak, ...) implements.

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
