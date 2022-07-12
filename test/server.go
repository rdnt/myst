package test

import (
	"fmt"

	"myst/src/server/application"
	"myst/src/server/repository"
	"myst/src/server/rest"
	"myst/src/server/rest/generated"
)

type Server struct {
	address string

	app    application.Application
	server *rest.Server
	client *generated.ClientWithResponses
}

func (s *IntegrationTestSuite) setupServer(address string) *Server {
	server, err := newServer(address)
	s.Require().Nil(err)

	return server
}

func newServer(address string) (*Server, error) {
	repo := repository.New()

	app, err := application.New(
		application.WithKeystoreRepository(repo),
		application.WithUserRepository(repo),
		application.WithInvitationRepository(repo),
	)
	if err != nil {
		return nil, err
	}

	server := rest.New(app)

	clientAddr := fmt.Sprintf("http://%s/api", address)

	client, err := generated.NewClientWithResponses(clientAddr)
	if err != nil {
		return nil, err
	}

	return &Server{
		address: address,
		app:     app,
		server:  server,
		client:  client,
	}, nil
}

func (s *Server) Start() error {
	return s.server.Start(s.address)
}

func (s *Server) Stop() error {
	return s.server.Stop()
}
