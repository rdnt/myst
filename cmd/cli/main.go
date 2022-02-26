package main

import (
	"github.com/spf13/cobra"

	"myst/internal/client/api/http/generated"
)

const address = "http://0.0.0.0:8081/api"

var client *generated.ClientWithResponses

func main() {
	var err error
	client, err = generated.NewClientWithResponses(address)
	if err != nil {
		panic(err)
	}

	rootCmd := &cobra.Command{
		Use: "myst",
	}

	rootCmd.AddCommand(authenticateCmd)
	rootCmd.AddCommand(keystoresCmd)

	err = rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
