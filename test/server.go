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

func (s *IntegrationTestSuite) setupServer() {
	keystoreRepo := keystorerepo.New()
	userRepo := userrepo.New()
	invitationRepo := invitationrepo.New()

	var err error
	s.server.app, err = application.New(
		application.WithKeystoreRepository(keystoreRepo),
		application.WithUserRepository(userRepo),
		application.WithInvitationRepository(invitationRepo),
	)
	s.Require().Nil(err)

	s.server.router = http.New(s.server.app).Engine
	s.server.server = httptest.NewServer(s.server.router)
}
