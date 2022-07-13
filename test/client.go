package test

import (
	"os"

	"github.com/gin-gonic/gin"

	"myst/src/client/application"
	"myst/src/client/remote"
	"myst/src/client/repository"
	"myst/src/client/rest"
	"myst/src/client/rest/generated"
)

type Client struct {
	address        string
	dir            string
	username       string
	password       string
	masterPassword string

	app     application.Application
	handler *gin.Engine
	server  *rest.Server
	client  *generated.ClientWithResponses
}

func (s *IntegrationTestSuite) setupClient(address string) *Client {
	username, password, masterPassword := s.rand(), s.rand(), s.rand()

	var err error
	dir, err := os.MkdirTemp("", "myst-client-data-*")
	s.Require().NoError(err)

	err = os.Chmod(dir, os.ModePerm)
	s.Require().NoError(err)

	s.T().Logf("Client starting (%s %s %s %s %s)...", address, dir, username, password, masterPassword)
	defer s.T().Logf("Client started (%s %s %s %s %s).", address, dir, username, password, masterPassword)

	client := &Client{
		dir:            dir,
		username:       username,
		password:       password,
		masterPassword: masterPassword,
	}

	repo, err := repository.New(client.dir)
	s.Require().NoError(err)

	rem, err := remote.New(
		remote.WithAddress("http://" + s.serverAddress + "/api"),
	)
	s.Require().NoError(err)

	client.app, err = application.New(
		application.WithKeystoreRepository(repo),
		application.WithRemote(rem),
	)
	s.Require().NoError(err)

	err = client.app.CreateEnclave(masterPassword)
	s.Require().NoError(err)

	client.restServer = rest.NewServer(client.app, nil)

	u, err := client.app.Register(username, password)
	s.Require().NoError(err)
	s.Require().Equal(username, u.Username)

	go func() {
		err = client.restServer.Start(address)
		s.Require().NoError(err)
	}()

	return client
}

func (s *IntegrationTestSuite) teardownClient(client *Client) {
	client.restServer.Stop()
	s.Require().NoError(os.RemoveAll(client.dir))
	// client.server.Close()
}
