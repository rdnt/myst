package test

import (
	"fmt"

	"github.com/gin-gonic/gin"

	application "myst/internal/server"
	"myst/internal/server/api/http"
	invitationrepo "myst/internal/server/core/invitationrepo/memory"
	keystorerepo "myst/internal/server/core/keystorerepo/memory"
	userrepo "myst/internal/server/core/userrepo/memory"
)

type Server struct {
	app     application.Application
	router  *gin.Engine
	address string
}

func (s *IntegrationTestSuite) setupServer(port int) *Server {
	server := &Server{}

	keystoreRepo := keystorerepo.New()
	userRepo := userrepo.New()
	invitationRepo := invitationrepo.New()

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
	//server.server.Close()
}
