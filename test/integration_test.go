package test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"myst/pkg/testing/capture"
)

func TestIntegration(t *testing.T) {
	s := &IntegrationTestSuite{
		capture: capture.New(t),
	}
	suite.Run(t, s)
}

func (s *IntegrationTestSuite) TestLogin() {
	s.T().Log("@@@@@@@@@@@@@@@@@")
	s.Require().Nil(nil)
}
