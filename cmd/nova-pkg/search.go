package main

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/J-V-Studios/nova-pkg/internal/backend"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <term>",
	Short: "Search all available backends",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		search(args[0])
	},
}

// result pairs one backend's search output with its error.
type result struct {
	pkgs []backend.Package
	err  error
	name string
}

// searchAll queries every available backend concurrently and returns
// the merged results sorted by backend registration order then name.
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
