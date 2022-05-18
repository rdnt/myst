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

func (s *IntegrationTestSuite) setupClient(serverAddress string) *Client {
	client := &Client{}

	var err error
	client.dir, err = os.MkdirTemp("", "myst-client-data-*")
	s.Require().Nil(err)

	err = os.Chmod(client.dir, os.ModePerm)
	s.Require().Nil(err)

	keystoreRepo, err := keystorerepo.New(client.dir)
	s.Require().Nil(err)

	keystoreService, err := keystoreservice.New(
		keystoreservice.WithKeystoreRepository(keystoreRepo),
	)
	s.Require().Nil(err)

	client.app, err = application.New(
		application.WithKeystoreService(keystoreService),
		application.WithRemoteAddress(serverAddress),
	)
	s.Require().Nil(err)

	client.router = clienthttp.New(client.app).Engine
	client.server = httptest.NewServer(s.server.router)

	return client
}

func (s *IntegrationTestSuite) teardownClient(client *Client) {
	s.Require().Nil(os.RemoveAll(client.dir))
	client.server.Close()
}
