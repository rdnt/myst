package test

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"myst/src/server/application"
	"myst/src/server/repository"
	"myst/src/server/rest"
)

type Server struct {
	app        application.Application
	router     *gin.Engine
	restServer *rest.Server
}

func (s *IntegrationTestSuite) setupServer(address string) *Server {
	s.T().Logf("Server starting (%s)...", address)
	defer s.T().Logf("Server started (%s).", address)

	server := &Server{}

	repo := repository.New()

	var err error
	server.app, err = application.New(
		application.WithKeystoreRepository(repo),
		application.WithUserRepository(repo),
		application.WithInvitationRepository(repo),
	)
	s.Require().NoError(err)

	server.restServer = rest.NewServer(server.app)

	go func() {
		fmt.Println(address)
		err = server.restServer.Start(address)
		s.Require().NoError(err)
	}()

	return server
}

func (s *IntegrationTestSuite) teardownServer(server *Server) {
	server.restServer.Stop()
}
