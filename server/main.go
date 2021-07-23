package main

import (
	"github.com/gin-gonic/gin"
	"myst/server/config"
	"myst/server/database"
	"myst/server/logger"
	"myst/server/regex"
	"myst/server/router"
	"net"
	"net/http"
)

func main() {
	err := config.Load()
	if err != nil {
		logger.Debugf("No environment files found. Using OS environment variables.")
	}

	err = logger.Init()
	if err != nil {
		logger.Errorf("Logger initialization failed: %s", err)
	}

	err = regex.Load()
	if err != nil {
		logger.Errorf("Regex initialization failed: %s", err)
		return
	}

	err = database.Init()
	if err != nil {
		logger.Errorf("Database initialization failed: %s", err)
	}

	mode := "release"
	if config.Debug == true {
		mode = "debug"
	}

	logger.Debugf("Starting server in %s mode...", mode)

	r := router.Init()
	err = Start(r)
	if err != nil {
		logger.Errorf("Failed to start server: %s", err)
		return
	}
}

func Start(r *gin.Engine) (err error) {
	// Get host and port from the environment
	var addr string

	netw := config.Get("SERVER_NETWORK")
	sock := config.Get("SERVER_SOCKET")
	host := config.Get("SERVER_HOST")
	port := config.Get("SERVER_HTTP_PORT")

	switch netw {
	case "unix":
		addr = sock
	default:
		netw = "tcp"
		if port == "" {
			port = "8625"
		}
		addr = host + ":" + port
	}
	ln, err := net.Listen(netw, addr)
	if err != nil {
		return err
	}
	logger.Debugf("Server started on port %s", port)

	err = http.Serve(ln, r)
	if err != nil {
		logger.Errorf("Failed to start HTTP server: %s", err)
	}

	return nil
}
