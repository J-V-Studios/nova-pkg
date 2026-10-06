package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed packages per backend",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		list()
	},
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
