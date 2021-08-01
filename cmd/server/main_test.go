package main

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"myst/cmd/server/api"
	"myst/pkg/config"
	"myst/pkg/router"
	"myst/pkg/testutil"
)

type IntegrationSuite struct {
	testutil.Suite
	server *httptest.Server
}

func TestIntegration(t *testing.T) {
	r := router.New(config.Debug)
	api.Init(r)

	s := testutil.NewSuite(
		testutil.WithName("Integration tests"),
		testutil.WithRouter(r),
	)
	time.Sleep(10 * time.Second)
	suite.Run(t, &IntegrationSuite{Suite: *s})
}

func (s *IntegrationSuite) SetupSuite() {
	s.Suite.SetupSuite()
	s.server = httptest.NewServer(s.Router())
}

func (s *IntegrationSuite) TearDownSuite() {
	s.server.Close()
	s.Suite.TearDownSuite()
}

func (s *IntegrationSuite) TestPing() {
	var res string
	s.GET(
		"/api/ping",
		nil,
		&res,
	)

	s.Require().Equal("Pong!", res)
	time.Sleep(1 * time.Second)
}

func (s *IntegrationSuite) TestPing2() {
	time.Sleep(1 * time.Second)
	var res string
	s.GET(
		"/api/ping",
		nil,
		&res,
	)

	s.Require().Equal("Pong!", res)
}

func (s *IntegrationSuite) TestPing3() {
	var res string
	s.GET(
		"/api/ping",
		nil,
		&res,
	)

	s.Require().Equal("Pong!", res)
}
