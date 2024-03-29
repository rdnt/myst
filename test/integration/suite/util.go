package suite

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"

	"myst/src/client/rest/generated"
)

func (s *Suite) CreateKeystore(t *testing.T) (keystore generated.Keystore) {
	keystoreName := strings.TrimSpace(s.Random(t))

	t.Run("Create a keystore", func(t *testing.T) {
		res, err := s.Client1.Client.CreateKeystoreWithResponse(s.Ctx,
			generated.CreateKeystoreJSONRequestBody{Name: keystoreName},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON201 != nil)
		assert.Equal(t, res.JSON201.Name, keystoreName)

		keystore = *res.JSON201
	})

	website, username, password, notes :=
		s.Random(t), s.Random(t), s.Random(t), s.Random(t)

	t.Run("Add an entry to the keystore", func(t *testing.T) {
		res, err := s.Client1.Client.CreateEntryWithResponse(s.Ctx, keystore.Id,
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
	})

	return keystore
}

func (s *Suite) CreateInvitation(t *testing.T, keystoreId string) (invitationId string) {
	t.Run("Invite someone to access the keystore", func(t *testing.T) {
		res, err := s.Client1.Client.CreateInvitationWithResponse(s.Ctx, keystoreId,
			generated.CreateInvitationJSONRequestBody{
				Invitee: s.Client2.Username,
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON201 != nil)

		invitationId = res.JSON201.Id
	})

	return invitationId
}
