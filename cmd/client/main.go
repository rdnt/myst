package main

import (
	"embed"
	"fmt"
	"os"
	"os/signal"

	"fyne.io/systray"
	"github.com/namsral/flag"
	"github.com/pkg/errors"

	"myst/pkg/logger"
	"myst/src/client/application"
	"myst/src/client/enclaverepo"
	"myst/src/client/remote"
	"myst/src/client/rest"
	"myst/src/client/scheduler"
)

//go:embed static/*
var static embed.FS

//go:embed icon.ico
var icon []byte

var log = logger.New("client", logger.Red)

type Config struct {
	RemoteAddress string
	Port          int
	DataDir       string
}

func parseFlags() Config {
	cfg := Config{}

	flag.StringVar(&cfg.RemoteAddress, "remote", "https://myst-abgx5.ondigitalocean.app/api", "URL address of the remote server API")
	flag.StringVar(&cfg.DataDir, "dir", ".",
		"Directory within which the user's enclave is stored. If the directory does not exist, it will be created.")
	flag.IntVar(&cfg.Port, "port", 8081, "Port the client should listen on")

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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	done := make(chan bool, 1)

	go systray.Run(onReady, func() {
		onExit(done)
	})

	select {
	case <-interrupt:
		// if process is interrupted, wait for systray to quit
		systray.Quit()
	case <-done:
	}

	err = cleanup()
	if err != nil {
		log.Error(err)
		os.Exit(1)
		return
	}
}

func onReady() {
	systray.SetIcon(icon)

	title := systray.AddMenuItem("Myst", "")
	title.Disable()

	systray.AddSeparator()

	quit := systray.AddMenuItem("Quit", "")

	go func() {
		<-quit.ClickedCh
		systray.Quit()
		return
	}()
}

func onExit(done chan bool) {
	done <- true
}

func run() (cleanup func() error, err error) {
	cfg := parseFlags()

	err = createDataDir(cfg.DataDir)
	if err != nil {
		return nil, errors.WithMessage(err, "unable to create data directory")
	}

	logger.EnableDebug = true

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

	sched, err := scheduler.New(app, rem)
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
