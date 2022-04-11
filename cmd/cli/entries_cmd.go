package main

import (
	"fmt"
	"myst/internal/client/api/http/generated"
	"net/http"

	"github.com/spf13/cobra"
)

var keystoreEntriesCmd = &cobra.Command{
	Use:     "keystoreEntries",
	Aliases: []string{"entries"},
	Run: func(cmd *cobra.Command, args []string) {

		res, err := client.CreateEntryWithResponse(
			cmd.Context(),
			args[0],
			generated.CreateEntryJSONRequestBody{
				Label:    args[1],
				Password: args[2],
				Username: args[3],
			},
		)

		if err != nil {
			panic(err)
		}

		if res.StatusCode() != http.StatusOK {
			fmt.Println("Error HTTP Status", res.StatusCode())
			return
		}

		fmt.Println("Entry added.")

	},
}
