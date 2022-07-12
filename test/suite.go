package test

import (
	"fmt"
	"testing"

	"github.com/phayes/freeport"
	"github.com/stretchr/testify/suite"

	"myst/pkg/config"
	"myst/pkg/crypto"
	"myst/pkg/logger"
	"myst/pkg/testing/capture"
	clientGenerated "myst/src/client/rest/generated"
	"myst/src/server/rest/generated"
)

type IntegrationTestSuite struct {
	suite.Suite
	capture *capture.Capture

	ports []int

	serverAddress  string
	client1address string
	client2address string
	client3address string

	// mini *miniredis.Miniredis

	_server  *Server
	_client1 *Client
	_client2 *Client
	_client3 *Client

	server  *generated.ClientWithResponses
	client1 *clientGenerated.ClientWithResponses
	client2 *clientGenerated.ClientWithResponses
	client3 *clientGenerated.ClientWithResponses
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

func (s *IntegrationTestSuite) rand(sizes ...int) string {
	size := 32
	if len(sizes) > 0 {
		size = sizes[0]
	}

	str, err := crypto.GenerateRandomString(size)
	s.Require().Nil(err)

	return str
}

func (s *IntegrationTestSuite) SetupSuite() {
	fmt.Println("Running integration_suite tests...")

	config.Debug = true
	logger.EnableDebug = config.Debug

	var err error
	s.ports, err = freeport.GetFreePorts(4)
	s.Require().NoError(err)

	s.serverAddress = fmt.Sprintf("localhost:%d", s.ports[0])
	s.client1address = fmt.Sprintf("localhost:%d", s.ports[1])
	s.client2address = fmt.Sprintf("localhost:%d", s.ports[2])
	s.client3address = fmt.Sprintf("localhost:%d", s.ports[3])

	s.server, err = generated.NewClientWithResponses("http://" + s.serverAddress + "/api")
	s.client1, err = clientGenerated.NewClientWithResponses("http://" + s.client1address + "/api")
	s.client2, err = clientGenerated.NewClientWithResponses("http://" + s.client2address + "/api")
	s.client3, err = clientGenerated.NewClientWithResponses("http://" + s.client3address + "/api")

	// s.mini, err = miniredis.Run()
	// s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	// s.mini.Close()

	// if verbose is enabled, print logger output
	if testing.Verbose() {
		fmt.Println()
		// fmt.Printf("\n%s\n", output)
	} else {
		fmt.Println()
	}
}

func (s *IntegrationTestSuite) SetupTest() {
	s.capture.Start()

	s._server = s.setupServer(s.serverAddress)
	s._client1 = s.setupClient(s.client1address)
	s._client2 = s.setupClient(s.client2address)
	s._client3 = s.setupClient(s.client3address)
}

func (s *IntegrationTestSuite) TearDownTest() {
	// start next tests with a flushed database
	// s.mini.FlushDB()
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

	s.teardownClient(s._client1)
	s.teardownClient(s._client2)
	s.teardownClient(s._client3)
	s.teardownServer(s._server)
}
