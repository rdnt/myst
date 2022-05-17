package test

import (
	"net/http/httptest"
	"os"

	clienthttp "myst/internal/client/api/http"
	"myst/internal/client/application"
	"myst/internal/client/application/keystoreservice"
	"myst/internal/client/keystorerepo"

	"github.com/gin-gonic/gin"
)

type Client struct {
	dir    string
	app    application.Application
	router *gin.Engine
	server *httptest.Server
}

func (s *IntegrationTestSuite) setupClient() {
	var err error
	s.client.dir, err = os.MkdirTemp("", "myst-client-data-*")
	s.Require().Nil(err)

	keystoreRepo, err := keystorerepo.New(s.client.dir)
	s.Require().Nil(err)

	keystoreService, err := keystoreservice.New(
		keystoreservice.WithKeystoreRepository(keystoreRepo),
	)
	s.Require().Nil(err)

	s.client.app, err = application.New(
		application.WithKeystoreService(keystoreService),
	)
	s.Require().Nil(err)

	s.client.router = clienthttp.New(s.client.app).Engine
	s.client.server = httptest.NewServer(s.server.router)
}

func (s *IntegrationTestSuite) teardownClient() {
	s.Require().Nil(os.RemoveAll(s.client.dir))
}
