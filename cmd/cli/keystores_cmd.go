package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var keystoresCmd = &cobra.Command{
	Use:     "keystores",
	Aliases: []string{"stores"},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: add flag to filter keystores

		res, err := client.KeystoresWithResponse(
			cmd.Context(),
		)
		if err != nil {
			panic(err)
		}

		if res.StatusCode() != http.StatusOK {
			fmt.Println("Error HTTP Status", res.StatusCode())
			return
		}

		b, err := json.MarshalIndent(*res.JSON200, "", "  ")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(b))
	},
}
