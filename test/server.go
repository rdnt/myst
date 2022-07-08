package test

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"myst/internal/server/application"
	"myst/internal/server/database/invitationrepo/memory"
	"myst/internal/server/database/keystorerepo/memory"
	"myst/internal/server/database/userrepo/memory"
)

type Server struct {
	app     application.Application
	router  *gin.Engine
	address string
}

func (s *IntegrationTestSuite) setupServer(port int) *Server {
	server := &Server{}

	keystoreRepo := keystorerepo.keystorerepo.New()
	userRepo := userrepo.userrepo.New()
	invitationRepo := invitationrepo.invitationrepo.New()

	var err error
	server.app, err = application.New(
		application.WithKeystoreRepository(keystoreRepo),
		application.WithUserRepository(userRepo),
		application.WithInvitationRepository(invitationRepo),
	)
	s.Require().NoError(err)

	api := http.New(server.app)
	server.address = fmt.Sprintf("localhost:%d", port)

	go func() {
		err = api.Run(server.address)
		s.Require().NoError(err)
	}()

	return server
}

func (s *IntegrationTestSuite) teardownServer(server *Server) {
	// server.server.Close()
}
