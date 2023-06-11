package main

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "myst",
}

func main() {
	rootCmd.AddCommand(bootstrapCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(lintCmd)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
