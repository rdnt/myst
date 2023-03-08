package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "myst",
}

func main() {
	rootCmd.AddCommand(bootstrapCmd)
	rootCmd.AddCommand(buildCmd)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done.")
}

// func Install{
// 	Step("Install", func() {
// 		CommandExists("go")
// 		Run("go mod download")
//
// 		CommandExists("npm")
// 		SetCurrentDir("ui")
// 		Run("npm ci")
// 		SetCurrentDir(".")
// 	})
// }
//
// func main() {
// 	defer Recover()
//
// 	Install()
//
// 	Step("Setup", func() {
// 		CleanDir("build")
// 		CleanDir("static")
// 		CleanDir("cmd/client/static")
// 	})
//
// 	if Env("GOOS") == "" {
// 		SetEnv("GOOS", "windows")
// 	}
//
// 	if Env("GOARCH") == "" {
// 		SetEnv("GOARCH", "amd64")
// 	}
//
// 	Step("Build server", func() {
// 		if runtime.GOOS == "windows" {
// 			Run("go build -o build/server/server-$GOOS-$GOARCH.exe cmd/server/main.go")
// 		} else {
// 			Run("go build -o build/server/server-$GOOS-$GOARCH cmd/server/main.go")
// 		}
// 	})
//
// 	Step("Build UI", func() {
// 		SetCurrentDir("ui")
// 		Run("npm run build")
// 		SetCurrentDir(".")
//
// 		CopyDir("static", "cmd/client/static")
// 	})
//
// 	Step("Build client", func() {
// 		if runtime.GOOS == "windows" {
// 			Run("go build -o build/client/client-$GOOS-$GOARCH.exe cmd/client/main.go")
// 		} else {
// 			Run("go build -o build/client/client-$GOOS-$GOARCH cmd/client/main.go")
// 		}
// 	})
// }
