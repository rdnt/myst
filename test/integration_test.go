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
	u1, err := s.server.app.CreateUser("rdnt", "1234")
	s.Require().NoError(err)

	err = s.client1.app.SignIn(u1.Id(), u1.Password())
	s.Require().NoError(err)

	u2, err := s.server.app.CreateUser("abcd", "5678")
	s.Require().NoError(err)

	err = s.client2.app.SignIn(u2.Id(), u2.Password())
	s.Require().NoError(err)

	k, err := s.client1.app.CreateFirstKeystore("test", "12345678")
	s.Require().NoError(err)
	s.Require().Equal(k.RemoteId(), "")

	err = s.client1.app.SignIn(u1.Id(), u1.Password())
	s.Require().NoError(err)

	inv, err := s.client1.app.CreateKeystoreInvitation(k.Id(), u2.Id())
	s.Require().NoError(err)

	s.T().Log(inv)

	// TODO: accept and finalize
}
