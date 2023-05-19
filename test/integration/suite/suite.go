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
	"myst/pkg/rand"
)

type Suite struct {
	Ctx     context.Context
	Server  *Server
	Client1 *Client
	Client2 *Client
	Client3 *Client
}

func New(t *testing.T) *Suite {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	config.Debug = true
	logger.EnableDebug = config.Debug

	ports, err := freeport.GetFreePorts(4)
	assert.NilError(t, err)

	addrs := lo.Map(ports, func(port int, _ int) string {
		return fmt.Sprintf("127.0.0.1:%d", port)
	})

	var server *Server
	var client1, client2, client3 *Client

	server = newServer(t, addrs[0])
	client1 = newClient(t, addrs[0], addrs[1])
	client2 = newClient(t, addrs[0], addrs[2])
	client3 = newClient(t, addrs[0], addrs[3])

	server.start(t)
	client1.start(t)
	client2.start(t)
	client3.start(t)

	t.Cleanup(func() {
		fmt.Println("cleanup")
		cancel()

		client1.stop(t)
		client2.stop(t)
		client3.stop(t)
		// stop the server after clients disconnect
		server.stop(t)
	})

	return &Suite{
		Ctx:     ctx,
		Server:  server,
		Client1: client1,
		Client2: client2,
		Client3: client3,
	}
}

// func (s *Suite) Run(t *testing.T, name string, fn func(*testing.T)) {
// 	if !t.Run(name, fn) {
// 		t.FailNow()
// 	}
// }

func (s *Suite) Random(t *testing.T) string {
	return random(t)
}

func random(t *testing.T) string {
	str, err := rand.String(16)
	assert.NilError(t, err)

	return str
}
