package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	logger2 "myst/pkg/logger"
)

var (
	httpSrv *http.Server
	log     = logger2.New("server", logger2.DefaultColor)

	ErrInvalidServerNetwork = fmt.Errorf("invalid server network")
)

func Start(r *gin.Engine) error {
	network := os.Getenv("MYST_SERVER_NETWORK_TYPE")
	if network == "" {
		network = "tcp"
	}

	var addr string
	var httpUrl *url.URL

	switch network {
	case "tcp", "tcp4", "tcp6":
		addr = os.Getenv("MYST_SERVER_ADDRESS")

		var err error
		httpUrl, err = ParseURL(addr)
		if err != nil {
			log.Error(err)
			return err
		}

		addr = fmt.Sprintf("%s:%s", httpUrl.Hostname(), httpUrl.Port())
	case "unix", "unixpacket":
		addr = os.Getenv("SERVER_SOCKET_PATH")
	default:
		log.Error(ErrInvalidServerNetwork)
		return ErrInvalidServerNetwork
	}

	var httpLn net.Listener

	var err error
	httpLn, err = net.Listen(network, addr)
	if err != nil {
		log.Error(err)
		return err
	}

	httpSrv = &http.Server{Handler: r}
	httpSrv.RegisterOnShutdown(
		func() {
			// cleanup upgraded/hijacked connections (e.g. websocket)
		},
	)

	switch network {
	case "tcp", "tcp4", "tcp6":
		log.Debugf(
			"Server started on port %s.\n",
			httpUrl.Port(),
		)
	case "unix", "unixpacket":
		log.Debugf("Server started.")
	}

	go func() {
		err := httpSrv.Serve(httpLn)
		if err != nil && err != http.ErrServerClosed {
			log.Error(err)
			Stop()
			return
		}
	}()

	return nil
}

func ParseURL(address string) (*url.URL, error) {
	scheme := "http://"
	if !strings.HasPrefix(address, scheme) {
		address = fmt.Sprintf("%s%s", scheme, address)
	}
	parsed, err := url.Parse(address)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if parsed.Port() == "" {
		address = fmt.Sprintf("%s%s:%d", scheme, parsed.Hostname(), 8080)
		parsed, err = url.Parse(address)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	return parsed, nil
}

func Stop() {
	// TODO wait at most 1 minute for server to shutdown
	err := httpSrv.Shutdown(context.Background())
	if err != nil {
		log.Error(err)
	}
}
