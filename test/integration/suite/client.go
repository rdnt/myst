package suite

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"testing"

	"gotest.tools/v3/assert"

	"myst/src/client/application"
	"myst/src/client/enclaverepo"
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

	SessionId []byte
}

func newClient(t *testing.T, serverAddress, address string) *Client {
	username, password, masterPassword := random(t), random(t), random(t)

	dir, err := os.MkdirTemp("", "myst-client-data-*")
	assert.NilError(t, err)

	err = os.Chmod(dir, os.ModePerm)
	assert.NilError(t, err)

	enclaveRepo := enclaverepo.New(dir)

	rem, err := remote.New(
		remote.WithAddress("http://" + serverAddress + "/api"),
	)
	assert.NilError(t, err)

	app := application.New(
		application.WithEnclave(enclaveRepo),
		application.WithRemote(rem),
	)

	server := rest.NewServer(app, nil)

	clientAddr := fmt.Sprintf("http://%s/api", address)

	c := &Client{
		dir: dir,

		Address:        address,
		Username:       username,
		Password:       password,
		MasterPassword: masterPassword,

		App:    app,
		Server: server,
	}

	client, err := generated.NewClientWithResponses(clientAddr, generated.WithRequestEditorFn(c.authenticationMiddleware()))
	assert.NilError(t, err)

	c.Client = client

	return c
}

func (c *Client) authenticationMiddleware() generated.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+base64.StdEncoding.EncodeToString(c.SessionId))
		return nil
	}
}

func (c *Client) start(t *testing.T) {
	err := c.Server.Start(c.Address)
	assert.NilError(t, err)
}

func (c *Client) stop(t *testing.T) {
	_ = c.Server.Stop()

	err := os.RemoveAll(c.dir)
	assert.NilError(t, err)
}
