package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	"myst/pkg/logger"
	"myst/pkg/rest"

	"github.com/mattn/go-colorable"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	name   string
	output *os.File
	input  *os.File
	router http.Handler
}

var (
	log = logger.New("test", logger.DefaultColor)
)

func NewSuite(opts ...func(*Suite)) *Suite {
	s := new(Suite)

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func WithName(name string) func(*Suite) {
	return func(s *Suite) {
		s.name = name
	}
}

func WithRouter(router http.Handler) func(*Suite) {
	return func(s *Suite) {
		s.router = router
	}
}

func (s *Suite) Router() http.Handler {
	return s.router
}

func (s *Suite) SetupSuite() {
	log.Printf("Running %s suite ...", s.name)

	enable := true
	colorable.EnableColorsStdout(&enable)
}

func (s *Suite) TearDownSuite() {}

func (s *Suite) SetupTest() {
	log.Printf("Running %s ...", s.T().Name())
}

func (s *Suite) TearDownTest() {
	var status string
	if s.T().Failed() {
		status = logger.Colorize("failed", logger.RedFg)
	} else if s.T().Skipped() {
		status = logger.Colorize("skipped", logger.RedFg)
	} else {
		status = logger.Colorize("passed", logger.GreenFg)
	}
	log.Printf("%s %s\n", s.T().Name(), status)
}

func (s *Suite) HandleStats(name string, stats *suite.SuiteInformation) {
	var status string
	if s.T().Failed() {
		status = logger.Colorize("failed", logger.RedFg)
	} else if s.T().Skipped() {
		status = logger.Colorize("skipped", logger.RedFg)
	} else {
		status = logger.Colorize("passed", logger.GreenFg)
	}

	log.Printf(
		"%s suite %s in %s",
		s.name,
		status,
		stats.End.Sub(stats.Start),
	)
}

func (s *Suite) GET(path string, body interface{}, dst interface{}) {
	s.Request("GET", path, body, dst)
}

func (s *Suite) Request(method string, path string, body interface{}, dst interface{}) {
	var b []byte

	if body != nil {
		var err error
		b, err = json.Marshal(body)
		s.Require().Equal(err, nil)
	}

	req, err := http.NewRequest(
		method,
		path,
		bytes.NewBuffer(b),
	)
	s.Require().Nil(err)

	resp := httptest.NewRecorder()
	s.router.ServeHTTP(resp, req)
	s.Require().Equal(http.StatusOK, resp.Code)

	s.ParseResponse(resp, &dst)
}

func (s *Suite) ParseResponse(resp *httptest.ResponseRecorder, r interface{}) {
	s.Require().Equal(http.StatusOK, resp.Code)

	var response rest.SuccessResponse
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	s.Require().Nil(err)
	s.Require().Equal("success", response.Status)

	b, err := json.Marshal(response.Data)
	s.Require().Nil(err)

	err = json.Unmarshal(b, &r)
	s.Require().Nil(err)
}
