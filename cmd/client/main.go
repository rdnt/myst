package main

import (
	"myst/internal/client"
	"myst/internal/client/api/http"
	keystorerepo "myst/internal/client/core/keystorerepo/fs"
	"myst/internal/client/core/keystoreservice"
	"myst/pkg/logger"
)

var log = logger.New("client", logger.Red)

func main() {
	keystoreRepo := keystorerepo.New()

	keystoreService, err := keystoreservice.New(
		keystoreservice.WithKeystoreRepository(keystoreRepo),
	)
	if err != nil {
		panic(err)
	}

	app, err := client.New(
		client.WithKeystoreRepository(keystoreRepo),
		client.WithKeystoreService(keystoreService),
	)
	if err != nil {
		panic(err)
	}

	api := http.New(app)

	err = api.Run(":8080")
	if err != nil {
		panic(err)
	}
}
