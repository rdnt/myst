package main

import (
	"github.com/spf13/cobra"

	. "myst/pkg/builder/util"
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Install dependencies and generate code",
	Run: func(cmd *cobra.Command, args []string) {
		CommandExists("go")
		CommandExists("npm")

		Step("Bootstrapping", func() {
			Run("go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest")

			Run("go mod download")
			Run("go generate ./...")

			Run("npm install -g openapi-typescript-codegen") // not working in WSL
			SetCurrentDir("ui")
			Run("npm ci")
			SetCurrentDir(".")
		})
	},
}
