package main

import (
	"fmt"
	"myst/internal/client/api/http/generated"
	"net/http"

	"github.com/spf13/cobra"
)

var authenticateCmd = &cobra.Command{
	Use:     "authenticate <password>",
	Aliases: []string{"auth"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: myst authenticate <password>")
			return
		}

		client, err := generated.NewClientWithResponses("http://0.0.0.0:8081/api")
		if err != nil {
			panic(err)
		}

		res, err := client.Authenticate(
			cmd.Context(),
			generated.AuthenticateJSONRequestBody{
				Password: args[0],
			},
		)
		if err != nil {
			panic(err)
		}

		if res.StatusCode != http.StatusOK {
			fmt.Println("Error HTTP Status", res.StatusCode)
			return
		}

		fmt.Println("Authenticated.")
	},
}
