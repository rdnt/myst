package main

import (
	"fmt"

	"myst/internal/client"
	keystorerepo "myst/internal/server/core/keystorerepo/memory"
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
