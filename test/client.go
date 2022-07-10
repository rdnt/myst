package test

import (
	"fmt"
	"os"

	"myst/src/client/application"
	"myst/src/client/remote"
	"myst/src/client/repository"
	"myst/src/client/rest"
)

type Client struct {
	dir            string
	app            application.Application
	address        string
	username       string
	password       string
	masterPassword string
}

func (s *IntegrationTestSuite) setupClient(serverAddress string, port int) *Client {
	username, password, masterPassword := s.rand(), s.rand(), s.rand()

	var err error
	dir, err := os.MkdirTemp("", "myst-client-data-*")
	s.Require().NoError(err)

	err = os.Chmod(dir, os.ModePerm)
	s.Require().NoError(err)

	s.T().Log("Client starting...", serverAddress, port, dir, username, password, masterPassword)
	defer s.T().Log("Client started", serverAddress, port, dir, username, password, masterPassword)

	client := &Client{
		dir:            dir,
		username:       username,
		password:       password,
		masterPassword: masterPassword,
	}

	// rem, err := remote.NewServer("http://localhost:8080")
	// s.Require().NoError(err)

	repo, err := repository.New(client.dir)
	s.Require().NoError(err)

	rem, err := remote.New(
		remote.WithAddress(serverAddress),
	)
	s.Require().NoError(err)

	client.app, err = application.New(
		application.WithKeystoreRepository(repo),
		application.WithRemote(rem),
	)
	s.Require().NoError(err)

	err = client.app.CreateEnclave(masterPassword)
	s.Require().NoError(err)

	srv := rest.NewServer(client.app, nil)
	client.address = fmt.Sprintf("localhost:%d", port)

	u, err := client.app.Register(username, password)
	s.Require().NoError(err)
	s.Require().Equal(username, u.Username)

	go func() {
		err = srv.Run(client.address)
		s.Require().NoError(err)
	}()

	return client
}

func (s *IntegrationTestSuite) teardownClient(client *Client) {
	s.Require().NoError(os.RemoveAll(client.dir))
	// client.server.Close()
}
