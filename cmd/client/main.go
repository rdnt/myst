package main

import (
	"embed"
	"fmt"
	"os"
	"os/signal"
	"time"

	"myst/pkg/config"
	"myst/pkg/logger"
	"myst/src/client/application"
	"myst/src/client/enclave"
	"myst/src/client/remote"
	"myst/src/client/rest"
	"myst/src/client/scheduler"

	"github.com/namsral/flag"
)

//go:embed static/*
var static embed.FS

type Config struct {
	RemoteAddress string
	Port          int
	DataDir       string
	Slow          bool
}

func parseFlags() Config {
	cfg := Config{}

	flag.StringVar(&cfg.RemoteAddress, "remote", "https://myst-abgx5.ondigitalocean.app/api", "URL address of the remote server API")
	flag.StringVar(&cfg.DataDir, "dir", "data", "Directory used to store the keystores")
	flag.IntVar(&cfg.Port, "port", 8080, "Port the client should listen on")
	flag.BoolVar(&cfg.Slow, "slow", false, "Wait 500ms before starting up")

	flag.Parse()

	return cfg
}

var log = logger.New("client", logger.Red)

func main() {
	cfg := parseFlags()

	if cfg.Slow {
		time.Sleep(500 * time.Millisecond)
	}

	logger.EnableDebug = config.Debug

	enc, err := enclave.New(cfg.DataDir)
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
		application.WithEnclave(enc),
		application.WithRemote(rem),
	)
	if err != nil {
		panic(err)
	}

	server := rest.NewServer(app, static)

	sched, err := scheduler.New(app)
	if err != nil {
		panic(err)
	}

	err = sched.Start()
	if err != nil {
		panic(err)
	}

	err = server.Start(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Error(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	_ = server.Stop()
	_ = sched.Stop()
}
