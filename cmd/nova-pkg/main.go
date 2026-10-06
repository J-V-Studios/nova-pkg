package main

import (
	"fmt"
	"os"
	"sort"
	"sync"

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
	fmt.Fprintln(os.Stderr, "  remove <name>   remove an installed package")
	fmt.Fprintln(os.Stderr, "  update          update all backends")
	fmt.Fprintln(os.Stderr, "  list            list installed packages")
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
	case "remove":
		if len(os.Args) != 3 {
			usage()
		}
		remove(os.Args[2])
	case "update":
		update()
	case "list":
		list()
	default:
		usage()
	}
}

// result pairs one backend's search output with its error.
type result struct {
	pkgs []backend.Package
	err  error
	name string
}

// searchAll queries every available backend concurrently and returns
// the merged results sorted by source then name.
func searchAll(query string) []result {
	var wg sync.WaitGroup
	ch := make(chan result, len(backends))
	for _, b := range backends {
		if !b.Available() {
			continue
		}
		wg.Add(1)
		go func(b backend.Backend) {
			defer wg.Done()
			pkgs, err := b.Search(query)
			ch <- result{pkgs: pkgs, err: err, name: b.Name()}
		}(b)
	}
	wg.Wait()
	close(ch)

	var results []result
	for r := range ch {
		results = append(results, r)
	}
	// Deterministic output: sort each backend's packages by name, and
	// order backends by the registration order in `backends`.
	order := make(map[string]int, len(backends))
	for i, b := range backends {
		order[b.Name()] = i
	}
	sort.Slice(results, func(i, j int) bool {
		return order[results[i].name] < order[results[j].name]
	})
	for _, r := range results {
		sort.Slice(r.pkgs, func(i, j int) bool {
			return r.pkgs[i].Name < r.pkgs[j].Name
		})
	}
	return results
}

func search(query string) {
	found := false
	for _, r := range searchAll(query) {
		if r.err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", r.name, r.err)
			continue
		}
		for _, p := range r.pkgs {
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

// update runs every available backend's update sequentially: sudo
// prompts must not race each other.
func update() {
	for _, b := range backends {
		if !b.Available() {
			continue
		}
		fmt.Printf("== %s ==\n", b.Name())
		if err := b.Update(); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", b.Name(), err)
		}
	}
}

func list() {
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
			fmt.Printf("%s/%s %s\n", p.Source, p.Name, p.Version)
		}
	}
}
