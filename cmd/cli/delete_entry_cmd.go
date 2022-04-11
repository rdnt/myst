package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var deleteEntryCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "rm"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Entry id required")
			return
		}

		entryId := args[0]

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

		if res.JSON200 == nil {
			fmt.Println("invalid response")
			return
		}

		var keystoreId string
		for _, k := range *res.JSON200 {
			for _, e := range k.Entries {
				if e.Id == entryId {
					keystoreId = k.Id
				}
			}
		}

		if keystoreId == "" {
			fmt.Println("Entry not found")
			return
		}

		res2, err := client.DeleteEntryWithResponse(
			cmd.Context(),
			keystoreId,
			entryId,
		)

		if err != nil {
			panic(err)
		}

		if res2.StatusCode() != http.StatusOK {
			fmt.Println("Error HTTP Status", res2.StatusCode())
			return
		}

		fmt.Println("Entry deleted.")
	},
}
