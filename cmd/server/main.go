package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/namsral/flag"

	"myst/pkg/config"
	"myst/pkg/logger"
	"myst/src/server/application"
	"myst/src/server/mongorepo"
	"myst/src/server/rest"
)

type Config struct {
	Port          int
	Slow          bool
	MongoAddress  string
	MongoDatabase string
}

func parseFlags() Config {
	cfg := Config{}

	flag.IntVar(&cfg.Port, "port", 8080, "Port the client should listen on")
	flag.BoolVar(&cfg.Slow, "slow", false, "Wait 500ms before starting up")

	flag.StringVar(&cfg.MongoAddress, "mongo-addr", "mongodb://localhost:27017",
		"The address of the mongodb server")
	flag.StringVar(&cfg.MongoDatabase, "mongo-db", "myst", "The name of the mongo database")

	flag.Parse()

	return cfg
}

var log = logger.New("app", logger.Red)

func main() {
	cfg := parseFlags()

	if cfg.Slow {
		time.Sleep(500 * time.Millisecond)
	}

	logger.EnableDebug = config.Debug

	repo, err := mongorepo.New(cfg.MongoAddress, cfg.MongoDatabase)
	if err != nil {
		panic(err)
	}

	app, err := application.New(
		application.WithKeystoreRepository(repo),
		application.WithUserRepository(repo),
		application.WithInvitationRepository(repo),
	)
	if err != nil {
		panic(err)
	}

	server := rest.NewServer(app)

	err = server.Run(fmt.Sprintf(":%d", cfg.Port))
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
