package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/labstack/gommon/log"
)

var (
	// httpSrv *http.Server
	// log     = logger.New("server", logger.DefaultColor)

	ErrInvalidNetwork = fmt.Errorf("invalid server network")
)

type Server struct {
	server   *http.Server
	listener net.Listener
}

func New(addr string, h http.Handler) (*Server, error) {
	httpLn, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	s := &Server{
		server:   &http.Server{Handler: h},
		listener: httpLn,
	}

	s.server.RegisterOnShutdown(
		func() {
			// cleanup upgraded/hijacked connections (e.g. websocket)
		},
	)

	go func() {
		err := s.server.Serve(httpLn)
		if err != nil && err != http.ErrServerClosed {
			log.Error(err)
			s.Stop()
			return
		}
	}()

	return s, nil
}

func (s *Server) Stop() error {
	// TODO wait at most 1 minute for server to shutdown
	err := s.server.Shutdown(context.Background())
	if err != nil {
		log.Error(err)
	}

	return s.listener.Close()
}
