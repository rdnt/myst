package integration_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/test/integration/suite"
)

func TestAuthorization(t *testing.T) {
	s := suite.New(t)

	err := s.Client1.App.Initialize(s.Client1.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client1.App.Register(s.Client1.Username, s.Client1.Password)
	assert.NilError(t, err)

	u, err := s.Client1.App.CurrentUser()
	assert.NilError(t, err)
	assert.Equal(t, u.Username, s.Client1.Username)

	err = s.Client1.App.SignOut()
	assert.NilError(t, err)

	u, err = s.Client1.App.CurrentUser()
	assert.NilError(t, err)
	assert.Assert(t, u == nil)

	_, err = s.Client1.App.SignIn(s.Client1.Username, s.Client1.Password)
	assert.NilError(t, err)

	u, err = s.Client1.App.CurrentUser()
	assert.NilError(t, err)
	assert.Equal(t, u.Username, s.Client1.Username)
}
