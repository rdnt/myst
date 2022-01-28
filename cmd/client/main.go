package main

import (
	"myst/internal/client/api/http"
	"myst/internal/client/application"
	"myst/internal/client/keystorerepo"
	"myst/pkg/logger"
)

var log = logger.New("client", logger.Red)

func main() {
	keystoreRepo, err := keystorerepo.New("data")
	if err != nil {
		panic(err)
	}

	app, err := application.New(
		application.WithKeystoreRepository(keystoreRepo),
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
