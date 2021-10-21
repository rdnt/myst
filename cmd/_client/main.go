package main

import (
	"fmt"

	keystorerepo "myst/app/server/core/keystorerepo/memory"

	"myst/app/client"
)

func main() {
	keystores := keystorerepo.New()

	fmt.Println(keystores.Keystores())

	fmt.Println(
		keystores.CreateKeystore(),
	)

	client.New()

	fmt.Println(keystores.Keystores())
}
