package main

import (
	"fmt"
	"myst/cmd/server/api"
	"myst/pkg/config"
	"myst/pkg/logger"
	"myst/pkg/router"
	"myst/pkg/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.Debugf("Starting server...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	err := logger.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := router.New(config.Debug)

	api.Init(r)

	err = server.Start(r)
	if err != nil {
		logger.Error(err)
		return
	}
	defer server.Stop()

	<-quit
	logger.Debugf("Server shutting down...")
}
