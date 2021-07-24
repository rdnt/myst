package main

import (
	"fmt"
	"myst/pkg/mongo"
	"os"
	"os/signal"
	"syscall"

	"myst/cmd/server/api"
	"myst/pkg/config"
	"myst/pkg/database"
	"myst/pkg/logger"
	"myst/pkg/regex"
	"myst/pkg/router"
	"myst/pkg/server"
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
	defer logger.Close()

	err = regex.Load()
	if err != nil {
		logger.Errorf("Regex initialization failed: %s", err)
		return
	}

	err = database.Init()
	if err != nil {
		logger.Errorf("Database initialization failed: %s", err)
		return
	}

	_, err = mongo.New("mongodb://localhost:27017")
	if err != nil {
		logger.Errorf("Database initialization failed: %s", err)
		return
	}
	defer mongo.Close()

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
