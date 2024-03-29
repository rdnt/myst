package integration_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/test/integration/suite"
)

func TestAuthentication(t *testing.T) {
	s := suite.New(t)

	sessionId, err := s.Client1.App.Initialize(s.Client1.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client1.App.Register(sessionId, s.Client1.Username, s.Client1.Password)
	assert.NilError(t, err)

	u, err := s.Client1.App.CurrentUser(sessionId)
	assert.NilError(t, err)
	assert.Equal(t, u.Username, s.Client1.Username)
}
