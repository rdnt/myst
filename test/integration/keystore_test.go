package integration_test

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"

	"myst/pkg/optional"
	"myst/src/client/rest/generated"
	"myst/test/integration/suite"
)

func TestKeystores(t *testing.T) {
	s := suite.New(t)

	err := s.Client1.App.Initialize(s.Client1.MasterPassword)
	assert.NilError(t, err)

	t.Run("Keystores are empty", func(t *testing.T) {
		res, err := s.Client1.Client.KeystoresWithResponse(s.Ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 0)
	})

	keystoreName := strings.TrimSpace(s.Random(t))
	var keystoreId string

	t.Run("Keystore created", func(t *testing.T) {
		res, err := s.Client1.Client.CreateKeystoreWithResponse(s.Ctx,
			generated.CreateKeystoreJSONRequestBody{Name: keystoreName},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON201 != nil)
		assert.Equal(t, res.JSON201.Name, keystoreName)

		keystoreId = (*res.JSON201).Id
	})

	t.Run("Keystore can be queried", func(t *testing.T) {
		res, err := s.Client1.Client.KeystoreWithResponse(s.Ctx, keystoreId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, res.JSON200.Id, keystoreId)
	})

	t.Run("Keystores contain created keystore", func(t *testing.T) {
		res, err := s.Client1.Client.KeystoresWithResponse(s.Ctx)
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

	t.Run("Entry created", func(t *testing.T) {
		res, err := s.Client1.Client.CreateEntryWithResponse(s.Ctx, keystoreId,
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

	t.Run("Entry can be queried", func(t *testing.T) {
		res, err := s.Client1.Client.KeystoreWithResponse(s.Ctx, keystoreId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(res.JSON200.Entries) == 1)
		assert.Equal(t, entryId, res.JSON200.Entries[0].Id)
	})

	newPassword := s.Random(t)
	newNotes := s.Random(t)

	t.Run("Entry can be updated", func(t *testing.T) {
		res, err := s.Client1.Client.UpdateEntryWithResponse(s.Ctx, keystoreId, entryId,
			generated.UpdateEntryJSONRequestBody{
				Password: optional.Ref(newPassword),
				Notes:    optional.Ref(newNotes),
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, newPassword, res.JSON200.Password)
		assert.Equal(t, newNotes, res.JSON200.Notes)

		res2, err := s.Client1.Client.KeystoreWithResponse(s.Ctx, keystoreId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Assert(t, len(res2.JSON200.Entries) == 1)
		assert.Equal(t, newPassword, res2.JSON200.Entries[0].Password)
		assert.Equal(t, newNotes, res2.JSON200.Entries[0].Notes)
	})

	t.Run("Entry can be deleted", func(t *testing.T) {
		_, err := s.Client1.Client.DeleteEntryWithResponse(s.Ctx, keystoreId, entryId)
		assert.NilError(t, err)

		res, err := s.Client1.Client.KeystoreWithResponse(s.Ctx, keystoreId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(res.JSON200.Entries) == 0)
	})

	t.Run("Keystore can be deleted", func(t *testing.T) {
		_, err := s.Client1.Client.DeleteKeystore(s.Ctx, keystoreId)
		assert.NilError(t, err)

		res, err := s.Client1.Client.KeystoresWithResponse(s.Ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 0)
	})
}
