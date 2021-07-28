package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mattn/go-colorable"
	"github.com/sanity-io/litter"

	"myst/cmd/server/api"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	"myst/config"
	"myst/logger"
	"myst/rest"
	"myst/router"
)

type IntegrationSuite struct {
	suite.Suite
	router   *gin.Engine
	server   *httptest.Server
	output   *os.File
	readPipe *os.File
}

func TestIntegration(t *testing.T) {
	s := new(IntegrationSuite)
	suite.Run(t, s)
}

func (s *IntegrationSuite) startCapture() {
	var err error
	s.readPipe, s.output, err = os.Pipe()
	s.Require().Nil(err)

	if testing.Verbose() {
		logger.StdoutWriter.SetWriter(s.output)
		logger.StderrWriter.SetWriter(s.output)
	} else {
		logger.StdoutWriter.SetWriter(nil)
		logger.StderrWriter.SetWriter(nil)
	}

	enable := true
	colorable.EnableColorsStdout(&enable)
}

func (s *IntegrationSuite) stopCapture() (string, error) {
	if s.output == nil || s.readPipe == nil {
		return "", errors.New("StartCapture not called before StopCapture")
	}
	err := s.output.Close()
	if err != nil {
		return "", err
	}
	b, err := ioutil.ReadAll(s.readPipe)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *IntegrationSuite) HandleStats(name string, stats *suite.SuiteInformation) {
	var status string
	litter.Dump(stats)
	if stats.Passed() {
		status = logger.Colorize("passed", logger.GreenFg)
	} else {
		status = logger.Colorize("failed", logger.RedFg)
	}

	fmt.Printf(
		"\n%s tests %s in %s\n\n",
		name,
		status,
		stats.End.Sub(stats.Start),
	)
}

func (s *IntegrationSuite) SetupSuite() {
	fmt.Println("Running integration tests...")
	if !testing.Verbose() {
		fmt.Println()
	}

	s.startCapture()

	s.router = router.New(config.Debug)

	api.Init(s.router)

	s.server = httptest.NewServer(s.router)
}

func (s *IntegrationSuite) TearDownSuite() {
	s.server.Close()
	logger.Close()

	output, err := s.stopCapture()
	s.Require().Nil(err)

	// if verbose is enabled, print logger output
	if testing.Verbose() {
		fmt.Printf("\n%s\n", output)
	} else {
		fmt.Println()
	}
}

func (s *IntegrationSuite) TearDownTest() {
	if !testing.Verbose() {
		// show progress
		s.T().Fail()
		if s.T().Failed() {
			fmt.Print(logger.Colorize("•", logger.RedFg))
		} else if s.T().Skipped() {
			fmt.Print(logger.Colorize("•", logger.WhiteFg))
		} else {
			fmt.Print(logger.Colorize("•", logger.GreenFg))
		}
	}
}

func (s *IntegrationSuite) TestPing() {
	panic("oops in tear down test")

	s.ping()
}

func (s *IntegrationSuite) TestFatal() {
	panic("test")
}

func (s *IntegrationSuite) ping() {
	req, err := http.NewRequest(
		"GET",
		"/api/ping",
		nil,
	)
	s.Require().Nil(err)

	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)
	s.Require().Equal(http.StatusOK, resp.Code)
}

func (s *IntegrationSuite) parseResponse(resp *httptest.ResponseRecorder, r interface{}) {
	var response rest.SuccessResponse
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	s.Require().Nil(err)
	s.Require().Equal("success", response.Status)

	b, err := json.Marshal(response.Data.(map[string]interface{}))
	s.Require().Nil(err)

	err = json.Unmarshal(b, &r)
	s.Require().Nil(err)
}
