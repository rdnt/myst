package main

import (
	"github.com/spf13/cobra"
)

var entryCmd = &cobra.Command{
	Use:     "entry",
	Aliases: []string{},
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// fmt.Println("entry called", args)

	// 	// if len(args) == 0 {
	// 	// 	panic("entry id is required")
	// 	// }
	// },
}
