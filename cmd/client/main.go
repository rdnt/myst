package main

import (
	"myst/internal/client/api/http"
	"myst/internal/client/application"
	"myst/internal/client/application/keystoreservice"
	"myst/internal/client/keystorerepo"
	"myst/pkg/config"
	"myst/pkg/logger"
)

var log = logger.New("client", logger.Red)

func main() {
	logger.EnableDebug = config.Debug

	keystoreRepo, err := keystorerepo.New("data")
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
		application.WithKeystoreService(keystoreService),
		application.WithRemoteAddress("http://localhost:8080"),
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
