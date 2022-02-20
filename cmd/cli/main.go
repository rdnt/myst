package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "myst",
	}

	rootCmd.AddCommand(authenticateCmd)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
