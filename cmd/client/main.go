package main

import (
	"embed"
	"fmt"

	"myst/pkg/config"
	"myst/pkg/logger"
	"myst/src/client/application"
	"myst/src/client/remote"
	"myst/src/client/repository"
	"myst/src/client/rest"

	"github.com/namsral/flag"
)

//go:embed static/*
var static embed.FS

type Config struct {
	RemoteAddress string
	Port          int
	DataDir       string
}

func parseFlags() Config {
	cfg := Config{}

	flag.StringVar(&cfg.RemoteAddress, "remote", "https://myst-abgx5.ondigitalocean.app/api", "URL address of the remote server API")
	flag.StringVar(&cfg.DataDir, "dir", "data", "Directory used to store the keystores")
	flag.IntVar(&cfg.Port, "port", 8080, "Port the client should listen on")

	flag.Parse()

	return cfg
}

var log = logger.New("client", logger.Red)

func main() {
	logger.EnableDebug = config.Debug

	cfg := parseFlags()

	// rem, err := remote.NewServer("http://localhost:8080")
	// if err != nil {
	//	panic(err)
	// }

	repo, err := repository.New(cfg.DataDir)
	if err != nil {
		panic(err)
	}

	rem, err := remote.New(
		remote.WithAddress(cfg.RemoteAddress),
	)
	if err != nil {
		panic(err)
	}

	app, err := application.New(
		application.WithKeystoreRepository(repo),
		application.WithRemote(rem),
	)
	if err != nil {
		panic(err)
	}

	server := rest.NewServer(app, static)

	err = server.Run(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		panic(err)
	}
}
