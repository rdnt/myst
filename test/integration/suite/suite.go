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
	"myst/src/client/rest/generated"
)

func random(t *testing.T) string {
	str, err := rand.String(16)
	assert.NilError(t, err)

	return str
}

type Suite struct {
	Server  *Server
	Client1 *Client
	Client2 *Client
	Client3 *Client
	Ctx     context.Context
}

func New(t *testing.T) *Suite {
	t.Parallel()

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
		cancel()

		client1.stop(t)
		client2.stop(t)
		client3.stop(t)
		// sto the server after clients disconnect
		server.stop(t)
	})

	return &Suite{
		Server:  server,
		Client1: client1,
		Client2: client2,
		Client3: client3,
		Ctx:     ctx,
	}
}

func (s *Suite) Run(t *testing.T, name string, fn func(*testing.T)) {
	if !t.Run(name, fn) {
		t.FailNow()
	}
}

func (s *Suite) Random(t *testing.T) string {
	return random(t)
}

func (s *Suite) CreateTestKeystore(t *testing.T) (keystoreId string) {
	keystoreName := s.Random(t)

	user := s.Client1

	s.Run(t, "Create a keystore", func(t *testing.T) {
		res, err := user.Client.CreateKeystoreWithResponse(s.Ctx,
			generated.CreateKeystoreJSONRequestBody{Name: keystoreName},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON201 != nil)
		assert.Equal(t, res.JSON201.Name, keystoreName)

		keystoreId = res.JSON201.Id
	})

	website, username, password, notes :=
		s.Random(t), s.Random(t), s.Random(t), s.Random(t)

	s.Run(t, "Add an entry to the keystore", func(t *testing.T) {
		res, err := user.Client.CreateEntryWithResponse(s.Ctx, keystoreId,
			generated.CreateEntryJSONRequestBody{
				Website:  website,
				Username: username,
				Password: password,
				Notes:    notes,
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON201 != nil)

		assert.Equal(t, res.JSON201.Website, website)
		assert.Equal(t, res.JSON201.Username, username)
		assert.Equal(t, res.JSON201.Password, password)
		assert.Equal(t, res.JSON201.Notes, notes)
	})

	return keystoreId
}
