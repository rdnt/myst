package main

import (
	"os"
	"os/signal"

	"myst/pkg/config"
	"myst/pkg/logger"
	"myst/src/server/application"
	"myst/src/server/repository"
	"myst/src/server/rest"
)

var log = logger.New("app", logger.Red)

func main() {
	logger.EnableDebug = config.Debug

	repo := repository.New()

	app, err := application.New(
		application.WithKeystoreRepository(repo),
		application.WithUserRepository(repo),
		application.WithInvitationRepository(repo),
	)
	if err != nil {
		panic(err)
	}

	server := rest.NewServer(app)

	err = server.Run(":8080")
	if err != nil {
		log.Error(err)
	}

	err = server.Start(":8080")
	if err != nil {
		log.Error(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	_ = server.Stop()
}
