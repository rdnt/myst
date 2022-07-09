package test

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"myst/src/server/application"
	"myst/src/server/repository"
	"myst/src/server/rest"
)

type Server struct {
	app     application.Application
	router  *gin.Engine
	address string
}

func (s *IntegrationTestSuite) setupServer(port int) *Server {
	server := &Server{}

	repo := repository.New()

	var err error
	server.app, err = application.New(
		application.WithKeystoreRepository(repo),
		application.WithUserRepository(repo),
		application.WithInvitationRepository(repo),
	)
	s.Require().NoError(err)

	srv := rest.NewServer(server.app)
	server.address = fmt.Sprintf("localhost:%d", port)

	go func() {
		err = srv.Run(server.address)
		s.Require().NoError(err)
	}()

	return server
}

func (s *IntegrationTestSuite) teardownServer(server *Server) {
	// server.server.Close()
}
