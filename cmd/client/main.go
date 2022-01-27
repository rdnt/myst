package main

import (
	application "myst/internal/client"
	"myst/internal/client/api/http"
	"myst/internal/client/core/enclaverepo"
	"myst/internal/client/core/keystoreservice"
	"myst/pkg/logger"
)

var log = logger.New("client", logger.Red)

func main() {
	keystoreRepo, err := enclaverepo.New("data")
	if err != nil {
		panic(err)
	}

	keystoreService, err := keystoreservice.New(
		keystoreservice.WithKeystoreRepository(keystoreRepo),
	)
	if err != nil {
		panic(err)
	}

	app, err := application.New(
		application.WithKeystoreRepository(keystoreRepo),
		application.WithKeystoreService(keystoreService),
	)
	if err != nil {
		panic(err)
	}

	api := http.New(app)

	err = api.Run(":8081")
	if err != nil {
		panic(err)
	}
}
