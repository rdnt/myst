package suite

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/phayes/freeport"
	"github.com/samber/lo"
	"gotest.tools/v3/assert"

	"myst/pkg/config"
	"myst/pkg/logger"
)

type Suite struct {
	Server  *Server
	Client1 *Client
	Client2 *Client
	Client3 *Client
	ctx     context.Context

	// Setup   func(*testing.T)
}

func setup(t *testing.T) Suite {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	config.Debug = true
	logger.EnableDebug = config.Debug

	ports, err := freeport.GetFreePorts(4)
	assert.NilError(t, err)

	addrs := lo.Map(ports, func(port int, _ int) string {
		return fmt.Sprintf("localhost:%d", port)
	})

	var server *Server
	var client1, client2, client3 *Client

	server, err = newServer(addrs[0])
	assert.NilError(t, err)

	client1, err = newClient(addrs[0], addrs[1])
	assert.NilError(t, err)

	client2, err = newClient(addrs[0], addrs[2])
	assert.NilError(t, err)

	client3, err = newClient(addrs[0], addrs[3])
	assert.NilError(t, err)

	// setup := func(t *testing.T) {
	err = server.Start()
	assert.NilError(t, err)

	err = client1.Start()
	assert.NilError(t, err)

	err = client2.Start()
	assert.NilError(t, err)

	err = client3.Start()
	assert.NilError(t, err)

	t.Cleanup(func() {
		now := time.Now()
		cancel()

		client1.Stop()
		client2.Stop()
		client3.Stop()
		server.Stop()

		fmt.Println("@@@@")
		fmt.Println(time.Since(now))
	})

	return Suite{
		Server:  server,
		Client1: client1,
		Client2: client2,
		Client3: client3,
		ctx:     ctx,
		// Setup:   setup,
	}

	// Step(t, "Test keystores", TestKeystores)
	// Step(t, "Test invitations", TestInvitations)
}

func startClient(t *testing.T, c *Client) {
	err := c.Start()
	assert.NilError(t, err)
}

func (s Suite) Run(t *testing.T, name string, fn func(*testing.T)) {
	if !t.Run(name, fn) {
		t.FailNow()
	}
}
