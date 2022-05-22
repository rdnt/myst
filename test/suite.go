package test

import (
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/suite"

	clientGenerated "myst/internal/client/api/http/generated"
	"myst/internal/server/api/http/generated"
	"myst/pkg/logger"
	"myst/pkg/testing/capture"
)

type IntegrationTestSuite struct {
	suite.Suite
	capture *capture.Capture

	mini *miniredis.Miniredis
	//router *gin.Engine
	//rdb      *redis.Client
	//server *httptest.Server

	_server  *Server
	_client1 *Client
	_client2 *Client

	server  *generated.ClientWithResponses
	client1 *clientGenerated.ClientWithResponses
	client2 *clientGenerated.ClientWithResponses
}

func (s *IntegrationTestSuite) HandleStats(name string, stats *suite.SuiteInformation) {
	var status string
	if stats.Passed() {
		status = logger.Colorize("passed", logger.Green)
	} else {
		status = logger.Colorize("failed", logger.Red)
	}

	fmt.Printf(
		"%s tests %s in %s\n",
		name,
		status,
		stats.End.Sub(stats.Start),
	)
}

func (s *IntegrationTestSuite) SetupSuite() {
	fmt.Println("Running integration_suite tests...")

	ports, err := freeport.GetFreePorts(3)
	s.Require().NoError(err)

	s._server = s.setupServer(ports[0])
	s._client1 = s.setupClient(s._server.server.URL, ports[1])
	s._client2 = s.setupClient(s._server.server.URL, ports[2])

	s.server, err = generated.NewClientWithResponses(s._server.server.URL)
	s.client1, err = clientGenerated.NewClientWithResponses(s._client1.server.URL)
	s.client2, err = clientGenerated.NewClientWithResponses(s._client2.server.URL)

	s.mini, err = miniredis.Run()
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	//s.server.Close()
	s.mini.Close()
	//logger.Close()

	s.teardownClient(s._client1)
	s.teardownClient(s._client2)
	s.teardownServer(s._server)

	//output := s.capture.Stop()

	// if verbose is enabled, print logger output
	if testing.Verbose() {
		fmt.Println()
		//fmt.Printf("\n%s\n", output)
	} else {
		fmt.Println()
	}
}

func (s *IntegrationTestSuite) SetupTest() {
	s.capture.Start()
}

func (s *IntegrationTestSuite) TearDownTest() {
	// start next tests with a flushed database
	s.mini.FlushDB()
	output := s.capture.Stop()

	if !testing.Verbose() {
		// show progress
		if s.T().Failed() {
			fmt.Print(logger.Colorize("•", logger.Red))
		} else if s.T().Skipped() {
			fmt.Print(logger.Colorize("•", logger.White))
		} else {
			fmt.Print(logger.Colorize("•", logger.Green))
		}
	} else {
		// if verbose is enabled, print logger output
		if testing.Verbose() {
			fmt.Printf("%s", output)
		} else {
			fmt.Println()
		}
	}
}
