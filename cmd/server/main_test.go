package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"myst/cmd/server/api"
	"myst/pkg/router"
	"myst/pkg/testutil"
)

type IntegrationSuite struct {
	testutil.Suite
}

func TestIntegration(t *testing.T) {
	r := router.New(
		router.WithDebug(testutil.Debug()),
	)

	api.Init(r)

	s := testutil.NewSuite(
		testutil.WithName("Integration tests"),
		testutil.WithRouter(r),
	)

	suite.Run(t, &IntegrationSuite{Suite: *s})
}

func (s *IntegrationSuite) SetupSuite() {
	s.Suite.SetupSuite()
}

func (s *IntegrationSuite) TearDownSuite() {
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
