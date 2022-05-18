package test

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"

	keystorerepo "myst/internal/server/core/keystorerepo/memory"

	application "myst/internal/server"
	"myst/internal/server/api/http"
	invitationrepo "myst/internal/server/core/invitationrepo/memory"
	userrepo "myst/internal/server/core/userrepo/memory"
)

type Server struct {
	app    application.Application
	router *gin.Engine
	server *httptest.Server
}

func (s *IntegrationTestSuite) setupServer() *Server {
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
	s.Require().Nil(err)

	server.router = http.New(server.app).Engine
	server.server = httptest.NewServer(server.router)

	return server
}

func (s *IntegrationTestSuite) teardownServer(server *Server) {
	server.server.Close()
}
