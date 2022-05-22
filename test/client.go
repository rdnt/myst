package test

import (
	"fmt"
	"os"

	"myst/internal/client/api/http"
	"myst/internal/client/application"
	"myst/internal/client/application/keystoreservice"
	"myst/internal/client/keystorerepo"
)

type Client struct {
	dir     string
	app     application.Application
	address string
}

func (s *IntegrationTestSuite) setupClient(serverAddress string, port int) *Client {
	client := &Client{}

	var err error
	client.dir, err = os.MkdirTemp("", "myst-client-data-*")
	s.Require().NoError(err)

	err = os.Chmod(client.dir, os.ModePerm)
	s.Require().NoError(err)

	//rem, err := remote.New("http://localhost:8080")
	//s.Require().NoError(err)

	keystoreRepo, err := keystorerepo.New(client.dir)
	s.Require().NoError(err)

	keystoreService, err := keystoreservice.New(
		keystoreservice.WithKeystoreRepository(keystoreRepo),
	)
	s.Require().NoError(err)

	client.app, err = application.New(
		application.WithKeystoreService(keystoreService),
		application.WithRemoteAddress(serverAddress),
	)
	s.Require().NoError(err)

	api := http.New(client.app)
	client.address = fmt.Sprintf("localhost:%d", port)

	go func() {
		err = api.Run(client.address)
		s.Require().NoError(err)
	}()

	return client
}

func (s *IntegrationTestSuite) teardownClient(client *Client) {
	s.Require().NoError(os.RemoveAll(client.dir))
	//client.server.Close()
}
