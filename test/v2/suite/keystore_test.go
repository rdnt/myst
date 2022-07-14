package suite

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/pkg/optional"
	"myst/src/client/rest/generated"
)

func TestKeystores(t *testing.T) {
	s := setup(t)

	user := s.Client1

	s.Run(t, "Keystores are empty", func(t *testing.T) {
		res, err := user.client.KeystoresWithResponse(s.ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 0)
	})

	keystoreName := random()
	var keystoreId string

	s.Run(t, "Keystore created", func(t *testing.T) {
		res, err := user.client.CreateKeystoreWithResponse(s.ctx,
			generated.CreateKeystoreJSONRequestBody{Name: keystoreName},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON201 != nil)
		assert.Equal(t, res.JSON201.Name, keystoreName)

		keystoreId = (*res.JSON201).Id
	})

	s.Run(t, "Keystore can be queried", func(t *testing.T) {
		res, err := user.client.KeystoreWithResponse(s.ctx, keystoreId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, res.JSON200.Id, keystoreId)
	})

	s.Run(t, "Keystores contain created keystore", func(t *testing.T) {
		res, err := user.client.KeystoresWithResponse(s.ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 1)
		assert.Equal(t, (*res.JSON200)[0].Id, keystoreId)
	})

	website := random()
	username := random()
	password := random()
	notes := random()
	var entryId string

	s.Run(t, "Entry created", func(t *testing.T) {
		res, err := user.client.CreateEntryWithResponse(s.ctx, keystoreId,
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

	s.Run(t, "Entry can be queried", func(t *testing.T) {
		res, err := user.client.KeystoreWithResponse(s.ctx, keystoreId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(res.JSON200.Entries) == 1)
		assert.Equal(t, entryId, res.JSON200.Entries[0].Id)
	})

	newPassword := random()
	newNotes := random()

	s.Run(t, "Entry can be updated", func(t *testing.T) {
		res, err := user.client.UpdateEntryWithResponse(s.ctx, keystoreId, entryId,
			generated.UpdateEntryJSONRequestBody{
				Password: optional.Ref(newPassword),
				Notes:    optional.Ref(newNotes),
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, newPassword, res.JSON200.Password)
		assert.Equal(t, newNotes, res.JSON200.Notes)

		res2, err := user.client.KeystoreWithResponse(s.ctx, keystoreId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Assert(t, len(res2.JSON200.Entries) == 1)
		assert.Equal(t, newPassword, res2.JSON200.Entries[0].Password)
		assert.Equal(t, newNotes, res2.JSON200.Entries[0].Notes)
	})

	s.Run(t, "Entry can be deleted", func(t *testing.T) {
		_, err := user.client.DeleteEntryWithResponse(s.ctx, keystoreId, entryId)
		assert.NilError(t, err)

		res, err := user.client.KeystoreWithResponse(s.ctx, keystoreId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(res.JSON200.Entries) == 0)
	})

	s.Run(t, "Keystore can be deleted", func(t *testing.T) {
		_, err := user.client.DeleteKeystore(s.ctx, keystoreId)
		assert.NilError(t, err)

		res, err := user.client.KeystoresWithResponse(s.ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 0)
	})
}
