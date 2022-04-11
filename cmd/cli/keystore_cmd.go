package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var keystoreCmd = &cobra.Command{
	Use:     "keystore",
	Aliases: []string{"store"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			panic("keystore command requires a keystore id")
		}

		res, err := client.KeystoreWithResponse(
			cmd.Context(),
			args[0],
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
