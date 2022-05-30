package main

import (
	"fmt"

	"myst/internal/client/api/http"
	"myst/internal/client/application"
	"myst/internal/client/application/keystoreservice"
	"myst/internal/client/keystorerepo"
	"myst/internal/client/remote"
	"myst/pkg/config"
	"myst/pkg/logger"

	"github.com/namsral/flag"
)

type Config struct {
	RemoteAddress  string
	RemoteUsername string
	RemotePassword string
	Port           int
	DataDir        string
}

func parseFlags() Config {
	cfg := Config{}

	flag.StringVar(&cfg.RemoteAddress, "remote", "http://localhost:8080", "URL address of the remote server")
	flag.IntVar(&cfg.Port, "port", 8081, "Port the client should listen on")

	flag.StringVar(&cfg.RemoteUsername, "username", "", "Username used to get authorized to the remote server")
	flag.StringVar(&cfg.RemotePassword, "password", "", "Password used to get authorized to the remote server")

	flag.StringVar(&cfg.DataDir, "dir", "data", "Directory used to store the keystores")

	flag.Parse()

	return cfg
}

var log = logger.New("client", logger.Red)

func main() {
	logger.EnableDebug = config.Debug

	cfg := parseFlags()

	//rem, err := remote.New("http://localhost:8080")
	//if err != nil {
	//	panic(err)
	//}

	keystoreRepo, err := keystorerepo.New(cfg.DataDir)
	if err != nil {
		panic(err)
	}

	keystoreService, err := keystoreservice.New(
		keystoreservice.WithKeystoreRepository(keystoreRepo),
	)
	if err != nil {
		panic(err)
	}

	remote, err := remote.New(
		remote.WithAddress(cfg.RemoteAddress),
		remote.WithUsername(cfg.RemoteUsername),
		remote.WithPassword(cfg.RemotePassword),
	)
	if err != nil {
		panic(err)
	}

	app, err := application.New(
		application.WithKeystoreService(keystoreService),
		application.WithRemote(remote),
	)
	if err != nil {
		panic(err)
	}

	api := http.New(app)

	err = api.Run(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		panic(err)
	}
}
