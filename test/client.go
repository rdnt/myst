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
	masterPassword string
}

func (s *IntegrationTestSuite) setupClient(serverAddress string, masterPassword, username, password string, port int) *Client {
	client := &Client{masterPassword: masterPassword}

	var err error
	client.dir, err = os.MkdirTemp("", "myst-client-data-*")
	s.Require().NoError(err)

	err = os.Chmod(client.dir, os.ModePerm)
	s.Require().NoError(err)

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

	err = client.app.SignIn(username, password)
	s.Require().NoError(err)

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
