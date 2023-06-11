package main

import (
	"github.com/spf13/cobra"
	"runtime"

	. "myst/pkg/builder/util"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build client and server",
	Run: func(cmd *cobra.Command, args []string) {
		CommandExists("go")
		CommandExists("npm")

		Step("Setting up", func() {
			CleanDir("build")
			CleanDir("static")
			CleanDir("cmd/client/static")
		})

		if Env("GOOS") == "" {
			SetEnv("GOOS", "windows")
		}

		if Env("GOARCH") == "" {
			SetEnv("GOARCH", "amd64")
		}

		Step("Building server", func() {
			if runtime.GOOS == "windows" {
				Run("go build -o build/server/server-$GOOS-$GOARCH.exe cmd/server/main.go")
			} else {
				Run("go build -o build/server/server-$GOOS-$GOARCH cmd/server/main.go")
			}
		})

		Step("Building UI", func() {
			SetCurrentDir("ui")
			Run("npm run build")
			SetCurrentDir(".")

			CopyDir("static", "cmd/client/static")

			// just to be sure the static dir will continue being pushed delete this file
			Touch("cmd/client/static/.gitkeep")
		})

		Step("Building client", func() {
			if runtime.GOOS == "windows" {
				Run("go build -o build/client/client-$GOOS-$GOARCH.exe cmd/client/main.go")
			} else {
				Run("go build -o build/client/client-$GOOS-$GOARCH cmd/client/main.go")
			}
		})
	},
}
