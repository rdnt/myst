package suite

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"

	"myst/pkg/uuid"
	"myst/src/server/application"
	"myst/src/server/mongorepo"
	"myst/src/server/rest"
	"myst/src/server/rest/generated"
)

type Server struct {
	Address string

	repo   *mongorepo.Repository
	App    application.Application
	Server *rest.Server
	Client *generated.ClientWithResponses
}

func newServer(t *testing.T, address string) *Server {
	repo, err := mongorepo.New("mongodb://localhost:27017", "myst-test-"+uuid.New().String())
	assert.NilError(t, err)

	app, err := application.New(
		application.WithKeystoreRepository(repo),
		application.WithUserRepository(repo),
		application.WithInvitationRepository(repo),
	)
	assert.NilError(t, err)

	server := rest.NewServer(app)

	clientAddr := fmt.Sprintf("http://%s/api", address)

	client, err := generated.NewClientWithResponses(clientAddr)
	assert.NilError(t, err)

	return &Server{
		Address: address,
		repo:    repo,
		App:     app,
		Server:  server,
		Client:  client,
	}
}

func (s *Server) start(t *testing.T) {
	err := s.Server.Start(s.Address)
	assert.NilError(t, err)

	err = s.repo.FlushDB()
	assert.NilError(t, err)
}

func (s *Server) stop(t *testing.T) {
	_ = s.Server.Stop()

	err := s.repo.FlushDB()
	assert.NilError(t, err)
}
