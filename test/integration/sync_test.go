package integration_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/src/client/rest/generated"
	"myst/test/integration/suite"
)

func TestKeystoreSyncUpStream(t *testing.T) {
	s := suite.New(t)

	err := s.Client1.App.Initialize(s.Client1.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client1.App.Register(s.Client1.Username, s.Client1.Password)
	assert.NilError(t, err)

	err = s.Client2.App.Initialize(s.Client2.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client2.App.Register(s.Client2.Username, s.Client2.Password)
	assert.NilError(t, err)

	var keystore generated.Keystore
	t.Run("Keystore is created", func(t *testing.T) {
		keystore = s.CreateKeystore(t)

		ks, err := s.Client1.App.Keystores()
		assert.NilError(t, err)
		assert.Assert(t, len(ks) == 1)
		assert.Equal(t, ks[keystore.Id].RemoteId, "")
		assert.Equal(t, ks[keystore.Id].Version, 2) // we added an entry
	})

	t.Run("Keystore is not uploaded yet", func(t *testing.T) {
		err := s.Client1.App.Sync()
		assert.NilError(t, err)

		rks, err := s.Server.App.Keystores()
		assert.NilError(t, err)
		assert.Equal(t, len(rks), 0)
	})

	t.Run("Create invitation for this keystore", func(t *testing.T) {
		_ = s.CreateInvitation(t, keystore.Id)
	})

	t.Run("Keystore is uploaded", func(t *testing.T) {
		rks, err := s.Server.App.Keystores()
		assert.NilError(t, err)
		assert.Equal(t, len(rks), 1)
		assert.Assert(t, rks[0].Id != "")
	})

	t.Run("Keystore is synced", func(t *testing.T) {
		err := s.Client1.App.Sync()
		assert.NilError(t, err)

		rks, err := s.Server.App.Keystores()
		assert.NilError(t, err)
		assert.Equal(t, len(rks), 1)
		assert.Assert(t, rks[0].Id != "")

		res, err := s.Client1.Client.KeystoresWithResponse(s.Ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 1)
		assert.Equal(t, (*res.JSON200)[0].RemoteId, rks[0].Id)

		ks, err := s.Client1.App.Keystores()
		assert.NilError(t, err)
		assert.Assert(t, len(ks) == 1)
		assert.Equal(t, ks[keystore.Id].Version, 3)
	})

	t.Run("Additional sync noop", func(t *testing.T) {
		err := s.Client1.App.Sync()
		assert.NilError(t, err)

		ks, err := s.Client1.App.Keystores()
		assert.NilError(t, err)
		assert.Assert(t, len(ks) == 1)
		assert.Equal(t, ks[keystore.Id].Version, 3)
	})
}
