package main

import (
	"github.com/spf13/cobra"
	. "myst/pkg/builder/util"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint go files",
	Run: func(cmd *cobra.Command, args []string) {
		CommandExists("go")
		CommandExists("golangci-lint")

		Step("Linting", func() {
			Run("golangci-lint run ./...")
		})
	},
}
