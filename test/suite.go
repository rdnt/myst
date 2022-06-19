package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/phayes/freeport"
	"github.com/stretchr/testify/suite"

	clientGenerated "myst/internal/client/api/http/generated"
	"myst/internal/server/api/http/generated"
	"myst/pkg/config"
	"myst/pkg/logger"
	"myst/pkg/testing/capture"
)

type IntegrationTestSuite struct {
	suite.Suite
	capture *capture.Capture

	//mini *miniredis.Miniredis

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

	config.Debug = true
	logger.EnableDebug = config.Debug

	ports, err := freeport.GetFreePorts(3)
	s.Require().NoError(err)

	s._server = s.setupServer(ports[0])
	s._client1 = s.setupClient("http://"+s._server.address, "rdnt", "1234", ports[1])
	s._client2 = s.setupClient("http://"+s._server.address, "abcd", "5678", ports[2])

	s.server, err = generated.NewClientWithResponses("http://" + s._server.address + "/api")
	s.client1, err = clientGenerated.NewClientWithResponses("http://" + s._client1.address + "/api")
	s.client2, err = clientGenerated.NewClientWithResponses("http://" + s._client2.address + "/api")

	//s.mini, err = miniredis.Run()
	//s.Require().NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, pass := "rdnttest", "1234"

	res, err := s.server.RegisterWithResponse(ctx, generated.RegisterJSONRequestBody{Username: user, Password: pass})
	s.Require().NoError(err)

	fmt.Println(string(res.Body))

	s.Require().NotNil(res.JSON201)
	s.Require().Equal(user, (*res.JSON201).Username)

	err = s._client1.app.SignIn()
	s.Require().NoError(err)

	user, pass = "abcdtest", "5678"

	res, err = s.server.RegisterWithResponse(ctx, generated.RegisterJSONRequestBody{Username: user, Password: pass})
	s.Require().NoError(err)

	s.Require().NotNil(res.JSON201)
	s.Require().Equal(user, (*res.JSON201).Username)

	err = s._client2.app.SignIn()
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	//s.mini.Close()

	s.teardownClient(s._client1)
	s.teardownClient(s._client2)
	s.teardownServer(s._server)

	// if verbose is enabled, print logger output
	if testing.Verbose() {
		fmt.Println()
		//fmt.Printf("\n%s\n", output)
	} else {
		fmt.Println()
	}
}

func (s *IntegrationTestSuite) SetupTest() {
	//s.capture.Start()

}

func (s *IntegrationTestSuite) TearDownTest() {
	// start next tests with a flushed database
	//s.mini.FlushDB()
	//output := s.capture.Stop()

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
			//fmt.Printf("%s", output)
		} else {
			fmt.Println()
		}
	}
}
