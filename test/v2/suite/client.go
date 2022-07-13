package suite

import (
	"fmt"
	"os"

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

	app    application.Application
	server *rest.Server
	client *generated.ClientWithResponses
}

func newClient(serverAddress, address string) (*Client, error) {
	strs, err := randomStrings(3, 32)
	if err != nil {
		return nil, err
	}

	username, password, masterPassword := strs[0], strs[1], strs[2]

	dir, err := os.MkdirTemp("", "myst-client-data-*")
	if err != nil {
		return nil, err
	}

	err = os.Chmod(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	repo, err := repository.New(dir)
	if err != nil {
		return nil, err
	}

	rem, err := remote.New(
		remote.WithAddress("http://" + serverAddress + "/api"),
	)
	if err != nil {
		return nil, err
	}

	app, err := application.New(
		application.WithKeystoreRepository(repo),
		application.WithRemote(rem),
	)
	if err != nil {
		return nil, err
	}

	server := rest.NewServer(app, nil)

	clientAddr := fmt.Sprintf("http://%s/api", address)

	client, err := generated.NewClientWithResponses(clientAddr)
	if err != nil {
		return nil, err
	}

	return &Client{
		address:        address,
		dir:            dir,
		username:       username,
		password:       password,
		masterPassword: masterPassword,

		app:    app,
		server: server,
		client: client,
	}, nil
}

func (c *Client) Start() error {
	return c.server.Start(c.address)
}

func (c *Client) Stop() error {
	return c.server.Stop()
}
