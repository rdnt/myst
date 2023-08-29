package main

import (
	"encoding/base64"
	"os"
	"os/signal"

	"github.com/namsral/flag"
	"github.com/pkg/errors"

	"myst/pkg/logger"
	"myst/src/server/application"
	"myst/src/server/mongorepo"
	"myst/src/server/rest"
)

var log = logger.New("app", logger.Red)

type Config struct {
	Port          int
	MongoAddress  string
	MongoDatabase string
	JWTSigningKey string
}

func parseFlags() Config {
	cfg := Config{}

	flag.IntVar(&cfg.Port, "port", 8080, "Port the client should listen on")

	flag.StringVar(&cfg.MongoAddress, "mongo-addr", "mongodb://localhost:27017",
		"The address of the mongodb server")
	flag.StringVar(&cfg.MongoDatabase, "mongo-db", "myst", "The name of the mongo database")

	flag.StringVar(&cfg.JWTSigningKey, "jwt-signing-key", "", "The key used for signing JWT tokens")

	flag.Parse()

	return cfg
}

func main() {
	cleanup, err := run()
	if err != nil {
		log.Error(err)
		os.Exit(1)
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	err = cleanup()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func run() (cleanup func() error, err error) {
	cfg := parseFlags()

	jwtSigningKey, err := base64.StdEncoding.DecodeString(cfg.JWTSigningKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode jwt signing key")
	}

	logger.EnableDebug = true

	repo, err := mongorepo.New(cfg.MongoAddress, cfg.MongoDatabase)
	if err != nil {
		return nil, errors.WithMessage(err, "could not create mongo repository")
	}

	app := application.New(
		application.WithKeystoreRepository(repo),
		application.WithUserRepository(repo),
		application.WithInvitationRepository(repo),
	)

	server := rest.NewServer(app, jwtSigningKey)

	err = server.Start(":8080")
	if err != nil {
		return nil, errors.WithMessage(err, "could not start server")
	}

	return func() error {
		err := server.Stop()
		if err != nil {
			return errors.WithMessage(err, "could not stop server")
		}

		return nil
	}, nil
}
