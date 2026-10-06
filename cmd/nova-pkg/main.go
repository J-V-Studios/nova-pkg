package main

import (
	"fmt"
	"os"

	"github.com/J-V-Studios/nova-pkg/internal/backend"
)

var backends = []backend.Backend{
	backend.Pacman{},
	backend.AUR{},
	backend.Flatpak{},
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: nova-pkg <command> [args]")
	fmt.Fprintln(os.Stderr, "  search <term>   search all available backends")
	fmt.Fprintln(os.Stderr, "  install <name>  install a package by name")
	os.Exit(2)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "search":
		if len(os.Args) != 3 {
			usage()
		}
		search(os.Args[2])
	case "install":
		if len(os.Args) != 3 {
			usage()
		}
		install(os.Args[2])
	default:
		usage()
	}
}

func search(query string) {
	found := false
	for _, b := range backends {
		if !b.Available() {
			continue
		}
		pkgs, err := b.Search(query)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", b.Name(), err)
			continue
		}
		for _, p := range pkgs {
			found = true
			fmt.Printf("%s/%s %s\n    %s\n", p.Source, p.Name, p.Version, p.Description)
		}
	}
	if !found {
		fmt.Println("no results")
	}
}

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
