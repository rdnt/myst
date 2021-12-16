package main

import (
	"myst/app/client"
	keystorerepo "myst/app/client/core/keystorerepo/fs"
	"myst/app/client/core/keystoreservice"
	"myst/app/client/rest"
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

	api := rest.New(app)

	err = api.Run(":8080")
	if err != nil {
		panic(err)
	}
}
