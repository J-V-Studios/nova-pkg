package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all backends",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		update()
	},
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
