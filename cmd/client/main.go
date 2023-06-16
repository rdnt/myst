package main

import (
	"embed"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/pkg/errors"

	"myst/pkg/config"
	"myst/pkg/logger"
	"myst/src/client/application"
	"myst/src/client/enclaverepo"
	"myst/src/client/remote"
	"myst/src/client/rest"
	"myst/src/client/scheduler"

	"github.com/namsral/flag"
)

//go:embed static/*
var static embed.FS

var log = logger.New("client", logger.Red)

type Config struct {
	RemoteAddress string
	Port          int
	DataDir       string
	Slow          bool
}

func parseFlags() Config {
	cfg := Config{}

	flag.StringVar(&cfg.RemoteAddress, "remote", "https://myst-abgx5.ondigitalocean.app/api", "URL address of the remote server API")
	flag.StringVar(&cfg.DataDir, "dir", ".",
		"Directory within which the user's enclave is stored. If the directory does not exist, it will be created.")
	flag.IntVar(&cfg.Port, "port", 8081, "Port the client should listen on")
	flag.BoolVar(&cfg.Slow, "slow", false, "Wait 500ms before starting up")

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
		return
	}
}

func run() (cleanup func() error, err error) {
	cfg := parseFlags()

	if cfg.Slow {
		time.Sleep(500 * time.Millisecond)
	}

	logger.EnableDebug = config.Debug

	err = createDataDir(cfg.DataDir)
	if err != nil {
		return nil, errors.WithMessage(err, "unable to create data directory")
	}

	enc := enclaverepo.New(cfg.DataDir)

	rem, err := remote.New(
		remote.WithAddress(cfg.RemoteAddress),
	)
	if err != nil {
		return nil, errors.WithMessage(err, "unable to create remote repository")
	}

	app := application.New(
		application.WithEnclave(enc),
		application.WithRemote(rem),
	)

	sched, err := scheduler.New(app)
	if err != nil {
		return nil, errors.WithMessage(err, "unable to create scheduler")
	}

	server := rest.NewServer(app, static)

	err = sched.Start()
	if err != nil {
		return nil, errors.WithMessage(err, "unable to start scheduler")
	}

	err = server.Start(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, errors.WithMessage(err, "unable to start server")
	}

	return func() error {
		err = server.Stop()
		if err != nil {
			return errors.WithMessage(err, "unable to stop server")
		}

		err = sched.Stop()
		if err != nil {
			return errors.WithMessage(err, "unable to stop scheduler")
		}

		return nil
	}, nil
}

func createDataDir(dir string) error {
	var create bool
	_, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		create = true
	} else if err != nil {
		return err
	}

	if !create {
		return nil
	}

	return os.Mkdir(dir, os.ModePerm)
}
