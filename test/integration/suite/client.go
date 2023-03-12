package suite

import (
	"fmt"
	"os"
	"testing"

	"gotest.tools/v3/assert"

	"myst/src/client/application"
	"myst/src/client/enclave"
	"myst/src/client/remote"
	"myst/src/client/rest"
	"myst/src/client/rest/generated"
)

type Client struct {
	dir string

	Address        string
	Username       string
	Password       string
	MasterPassword string

	App    application.Application
	Server *rest.Server
	Client *generated.ClientWithResponses
}

func newClient(t *testing.T, serverAddress, address string) *Client {
	username, password, masterPassword := random(t), random(t), random(t)

	dir, err := os.MkdirTemp("", "myst-client-data-*")
	assert.NilError(t, err)

	err = os.Chmod(dir, os.ModePerm)
	assert.NilError(t, err)

	repo, err := enclave.New(dir)
	assert.NilError(t, err)

	rem, err := remote.New(
		remote.WithAddress("http://" + serverAddress + "/api"),
	)
	assert.NilError(t, err)

	app, err := application.New(
		application.WithEnclave(repo),
		application.WithRemote(rem),
		application.WithInvitationRepository(rem),
		application.WithRepository(repo),
		application.WithCredentials(repo),
	)
	assert.NilError(t, err)

	server := rest.NewServer(app, nil)

	clientAddr := fmt.Sprintf("http://%s/api", address)

	client, err := generated.NewClientWithResponses(clientAddr)
	assert.NilError(t, err)

	return &Client{
		dir: dir,

		Address:        address,
		Username:       username,
		Password:       password,
		MasterPassword: masterPassword,

		App:    app,
		Server: server,
		Client: client,
	}
}

func (c *Client) start(t *testing.T) {
	err := c.Server.Start(c.Address)
	assert.NilError(t, err)

	err = c.App.Initialize(c.MasterPassword)
	assert.NilError(t, err)

	_, err = c.App.Register(c.Username, c.Password)
	assert.NilError(t, err)
}

func (c *Client) stop(t *testing.T) {
	_ = c.Server.Stop()

	err := os.RemoveAll(c.dir)
	assert.NilError(t, err)
}
