package test

import (
	"fmt"
	"net/http/httptest"

	"github.com/gin-gonic/gin"

	application "myst/internal/server"
	"myst/internal/server/api/http"
	invitationrepo "myst/internal/server/core/invitationrepo/memory"
	keystorerepo "myst/internal/server/core/keystorerepo/memory"
	userrepo "myst/internal/server/core/userrepo/memory"
)

type Server struct {
	app    application.Application
	router *gin.Engine
	server *httptest.Server
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

	go func() {
		err = api.Run(fmt.Sprintf(":%d", port))
		s.Require().NoError(err)
	}()

	return server
}

func (s *IntegrationTestSuite) teardownServer(server *Server) {
	server.server.Close()
}
