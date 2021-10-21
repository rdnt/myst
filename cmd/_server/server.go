package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"myst/pkg/mongo"

	"github.com/gin-gonic/gin"

	"myst/cmd/_server/api"
	"myst/pkg/config"
	"myst/pkg/logger"
	"myst/pkg/regex"
	"myst/pkg/router"
	"myst/pkg/server"
	"myst/pkg/storage"
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

	err = storage.Init()
	if err != nil {
		logger.Errorf("Storage initialization failed: %s", err)
		return
	}

	db, err := mongo.New(mongo.WithURI("mongodb://localhost:27017"))
	if err != nil {
		logger.Errorf("Database initialization failed: %s", err)
		return
	}
	defer db.Close()

	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := router.New(router.WithDebug(config.Debug))

	api.Init(r)

	err = server.Start(r)
	if err != nil {
		return
	}
	defer server.Stop()

	<-quit
	logger.Debugf("Server shutting down...")
}
