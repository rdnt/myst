package main

import (
	"fmt"

	"myst/app/client"

	keystorerepo "myst/app/server/keystorerepo/memory"
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
