package main

import "github.com/spf13/cobra"

var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove an installed package",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		remove(args[0])
	},
}
