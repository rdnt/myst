package suite

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/pkg/optional"
	"myst/src/client/rest/generated"
	"myst/test/integration/suite"
)

func TestKeystores(t *testing.T) {
	s := suite.New(t)

	user := s.Client1

	s.Run("Keystores are empty", func(t *testing.T) {
		res, err := user.Client.KeystoresWithResponse(s.Ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 0)
	})

	keystoreName := s.Random(t)
	var keystoreId string

	s.Run("Keystore created", func(t *testing.T) {
		res, err := user.Client.CreateKeystoreWithResponse(s.Ctx,
			generated.CreateKeystoreJSONRequestBody{Name: keystoreName},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON201 != nil)
		assert.Equal(t, res.JSON201.Name, keystoreName)

		keystoreId = (*res.JSON201).Id
	})

	s.Run("Keystore can be queried", func(t *testing.T) {
		res, err := user.Client.KeystoreWithResponse(s.Ctx, keystoreId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, res.JSON200.Id, keystoreId)
	})

	s.Run("Keystores contain created keystore", func(t *testing.T) {
		res, err := user.Client.KeystoresWithResponse(s.Ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 1)
		assert.Equal(t, (*res.JSON200)[0].Id, keystoreId)
	})

	website := s.Random(t)
	username := s.Random(t)
	password := s.Random(t)
	notes := s.Random(t)
	var entryId string

	s.Run("Entry created", func(t *testing.T) {
		res, err := user.Client.CreateEntryWithResponse(s.Ctx, keystoreId,
			generated.CreateEntryJSONRequestBody{
				Website:  website,
				Username: username,
				Password: password,
				Notes:    notes,
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON201 != nil)
		assert.Equal(t, res.JSON201.Website, website)
		assert.Equal(t, res.JSON201.Username, username)
		assert.Equal(t, res.JSON201.Password, password)
		assert.Equal(t, res.JSON201.Notes, notes)

		entryId = (*res.JSON201).Id
	})

	s.Run("Entry can be queried", func(t *testing.T) {
		res, err := user.Client.KeystoreWithResponse(s.Ctx, keystoreId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(res.JSON200.Entries) == 1)
		assert.Equal(t, entryId, res.JSON200.Entries[0].Id)
	})

	newPassword := s.Random(t)
	newNotes := s.Random(t)

	s.Run("Entry can be updated", func(t *testing.T) {
		res, err := user.Client.UpdateEntryWithResponse(s.Ctx, keystoreId, entryId,
			generated.UpdateEntryJSONRequestBody{
				Password: optional.Ref(newPassword),
				Notes:    optional.Ref(newNotes),
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, newPassword, res.JSON200.Password)
		assert.Equal(t, newNotes, res.JSON200.Notes)

		res2, err := user.Client.KeystoreWithResponse(s.Ctx, keystoreId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Assert(t, len(res2.JSON200.Entries) == 1)
		assert.Equal(t, newPassword, res2.JSON200.Entries[0].Password)
		assert.Equal(t, newNotes, res2.JSON200.Entries[0].Notes)
	})

	s.Run("Entry can be deleted", func(t *testing.T) {
		_, err := user.Client.DeleteEntryWithResponse(s.Ctx, keystoreId, entryId)
		assert.NilError(t, err)

		res, err := user.Client.KeystoreWithResponse(s.Ctx, keystoreId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(res.JSON200.Entries) == 0)
	})

	s.Run("Keystore can be deleted", func(t *testing.T) {
		_, err := user.Client.DeleteKeystore(s.Ctx, keystoreId)
		assert.NilError(t, err)

		res, err := user.Client.KeystoresWithResponse(s.Ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 0)
	})
}
