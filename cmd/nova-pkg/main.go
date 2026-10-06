package main

import (
	"os"

	"github.com/J-V-Studios/nova-pkg/internal/backend"
	"github.com/spf13/cobra"
)

var backends = []backend.Backend{
	backend.Pacman{},
	backend.AUR{},
	backend.Flatpak{},
}

var rootCmd = &cobra.Command{
	Use:   "nova-pkg",
	Short: "Universal package installer for Arch Linux",
	Long:  "nova-pkg wraps pacman, AUR and Flatpak behind one CLI.",
}

func init() {
	rootCmd.AddCommand(searchCmd, installCmd, removeCmd, updateCmd, listCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
