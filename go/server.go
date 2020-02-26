package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sht/myst/go/config"
	"github.com/sht/myst/go/logger"
	"github.com/sht/myst/go/regex"
	"github.com/sht/myst/go/router"
	"net"
	"net/http"
)

func main() {

	err := config.Load()
	logger.Init()
	if err != nil {
		logger.Logf("STARTUP", "No environment files found. Using OS environment variables.")
	}
	err = regex.Load()
	if err != nil {
		logger.Errorf("STARTUP", "Regex validation failed: %s", err)
		return
	}

	mode := "release"
	if config.Debug == true {
		mode = "debug"
	}
	logger.Logf("STARTUP", "Starting server in %s mode...", mode)

	r := router.Init()
	err = Start(r)
	if err != nil {
		logger.Errorf("FATAL", "Failed to start server: %s", err)
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
	logger.Logf("STARTUP", "Server started on port %s", port)

	err = http.Serve(ln, r)
	if err != nil {
		logger.Errorf("STARTUP", "Failed to start HTTP server: %s", err)
	}

	return nil
}
