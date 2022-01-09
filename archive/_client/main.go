package main

import (
	"fmt"

	application "myst/internal/client"
	keystorerepo "myst/internal/server/core/keystorerepo/memory"
)

func main() {
	keystores := keystorerepo.New()

	fmt.Println(keystores.Keystores())

	fmt.Println(
		keystores.CreateKeystore(),
	)

	application.New()

	fmt.Println(keystores.Keystores())
}
