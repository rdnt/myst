package suite

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"

	"myst/src/server/application"
	"myst/src/server/repository"
	"myst/src/server/rest"
	"myst/src/server/rest/generated"
)

type Server struct {
	Address string

	App    application.Application
	Server *rest.Server
	Client *generated.ClientWithResponses
}

func newServer(t *testing.T, address string) *Server {
	repo := repository.New()

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

		App:    app,
		Server: server,
		Client: client,
	}
}

func (s *Server) start(t *testing.T) {
	err := s.Server.Start(s.Address)
	assert.NilError(t, err)
}

func (s *Server) stop(t *testing.T) {
	_ = s.Server.Stop()
}
