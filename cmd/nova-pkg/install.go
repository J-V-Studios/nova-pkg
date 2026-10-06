package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install <name>",
	Short: "Install a package by exact name",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		install(args[0])
	},
}

// install finds the first backend (registration order) with an exact
// name match and installs from there.
func install(name string) {
	for _, b := range backends {
		if !b.Available() {
			continue
		}
		pkgs, err := b.Search(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", b.Name(), err)
			continue
		}
		for _, p := range pkgs {
			if p.Name != name {
				continue
			}
			fmt.Printf("installing %s from %s\n", p.Name, p.Source)
			if err := b.Install(name); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
	}
	fmt.Fprintf(os.Stderr, "package %q not found in any backend\n", name)
	os.Exit(1)
}

// remove finds the backend where name is installed (exact match,
// registration order) and removes it from there.
func remove(name string) {
	for _, b := range backends {
		if !b.Available() {
			continue
		}
		pkgs, err := b.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", b.Name(), err)
			continue
		}
		for _, p := range pkgs {
			if p.Name != name {
				continue
			}
			fmt.Printf("removing %s from %s\n", p.Name, p.Source)
			if err := b.Remove(name); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
	}
	fmt.Fprintf(os.Stderr, "package %q is not installed\n", name)
	os.Exit(1)
}
