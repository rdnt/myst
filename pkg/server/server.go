package server

import (
	"context"
	"net"
	"net/http"

	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

type Server struct {
	server   *http.Server
	listener net.Listener
}

func New(addr string, h http.Handler) (*Server, error) {
	httpLn, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to listen to address: %s", addr)
	}

	s := &Server{
		server:   &http.Server{Handler: h},
		listener: httpLn,
	}

	// s.server.RegisterOnShutdown(
	// 	func() {
	// 		// cleanup upgraded/hijacked connections (e.g. websocket)
	// 	},
	// )

	go func() {
		err := s.server.Serve(httpLn)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error(err)
			_ = s.Stop()
			return
		}
	}()

	return s, nil
}

func (s *Server) Stop() error {
	err := s.server.Shutdown(context.Background())
	if err != nil {
		return err
	}

	return nil
}
