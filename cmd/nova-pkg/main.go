package main

import (
	"fmt"
	"os"

	"github.com/J-V-Studios/nova-pkg/internal/backend"
	"github.com/J-V-Studios/nova-pkg/internal/config"
	"github.com/spf13/cobra"
)

// defaultBackends is the built-in order; config.Priority reorders it.
var defaultBackends = []backend.Backend{
	backend.Pacman{},
	backend.AUR{},
	backend.Flatpak{},
}

var backends = defaultBackends

var rootCmd = &cobra.Command{
	Use:   "nova-pkg",
	Short: "Universal package installer for Arch Linux",
	Long:  "nova-pkg wraps pacman, AUR and Flatpak behind one CLI.",
}

// orderBackends returns defaultBackends reordered so backends named in
// priority come first, in that order. Unknown names are ignored;
// unlisted backends keep their relative default order at the end.
func orderBackends(priority []string) []backend.Backend {
	if len(priority) == 0 {
		return defaultBackends
	}
	byName := make(map[string]backend.Backend, len(defaultBackends))
	for _, b := range defaultBackends {
		byName[b.Name()] = b
	}
	var ordered []backend.Backend
	used := make(map[string]bool, len(priority))
	for _, name := range priority {
		if b, ok := byName[name]; ok && !used[name] {
			ordered = append(ordered, b)
			used[name] = true
		}
	}
	for _, b := range defaultBackends {
		if !used[b.Name()] {
			ordered = append(ordered, b)
		}
	}
	return ordered
}

func init() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}
	backends = orderBackends(cfg.Priority)
	rootCmd.AddCommand(searchCmd, installCmd, removeCmd, updateCmd, listCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
