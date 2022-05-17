package test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	"myst/pkg/logger"
	"myst/test/pkg/capture"
)

type IntegrationTestSuite struct {
	suite.Suite
	capture *capture.Capture

	mini   *miniredis.Miniredis
	router *gin.Engine
	//rdb      *redis.Client
	server *httptest.Server
}

func TestIntegration(t *testing.T) {
	s := &IntegrationTestSuite{}

	s.capture = capture.New(t)

	suite.Run(t, s)
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
	fmt.Println("Running integration tests...")

	//err := logger.Init()
	//s.Require().Nil(err)

	//s.capture.Start()

	var err error
	s.mini, err = miniredis.Run()
	s.Require().Nil(err)

	//s.rdb = redis.NewClient(&redis.Options{
	//	Network: "tcp",
	//	Addr:    s.mini.Addr(),
	//})

	//database.SetDB(s.rdb)
	//models.Init(s.rdb)

	//s.router = router.Init()
	//s.server = httptest.NewServer(s.router)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	//s.server.Close()
	s.mini.Close()
	//logger.Close()

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

func (s *IntegrationTestSuite) TestLogin() {
	fmt.Print("@@@@@@@@@@@@@@@@@")
	s.Require().Nil(nil)
}
