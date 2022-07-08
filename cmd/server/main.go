package main

import (
	"myst/internal/server/application"
	"myst/internal/server/memdb"
	"myst/internal/server/rest"
	"myst/pkg/config"
	"myst/pkg/logger"
)

var log = logger.New("app", logger.Red)

func main() {
	logger.EnableDebug = config.Debug

	mem := memdb.New()

	app, err := application.New(
		application.WithKeystoreRepository(mem),
		application.WithUserRepository(mem),
		application.WithInvitationRepository(mem),
	)
	if err != nil {
		panic(err)
	}

	server := rest.NewServer(app)

	err = server.Run(":8080")
	if err != nil {
		panic(err)
	}
}
