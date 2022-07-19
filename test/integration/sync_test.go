package integration_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/src/client/rest/generated"
)

func TestKeystoreSyncUpStream(t *testing.T) {
	s := setup(t)

	var keystore generated.Keystore
	s.Run(t, "Keystore is created", func(t *testing.T) {
		keystore = s.createKeystore(t)

		ks, err := s.Client1.App.Keystores()
		assert.NilError(t, err)
		assert.Assert(t, len(ks) == 1)
		assert.Equal(t, ks[keystore.Id].RemoteId, "")
		assert.Equal(t, ks[keystore.Id].Version, 2) // we added an entry
	})

	s.Run(t, "Keystore is uploaded", func(t *testing.T) {
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

	s.Run(t, "Additional sync noop", func(t *testing.T) {
		err := s.Client1.App.Sync()
		assert.NilError(t, err)

		ks, err := s.Client1.App.Keystores()
		assert.NilError(t, err)
		assert.Assert(t, len(ks) == 1)
		assert.Equal(t, ks[keystore.Id].Version, 3)
	})
}
