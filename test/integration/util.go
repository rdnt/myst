package suite

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/src/client/rest/generated"
	"myst/test/integration/suite"
)

func CreateKeystore(s *suite.Suite) (keystoreId string) {
	keystoreName := s.Random(s.T)

	user := s.Client1

	s.Run(s.T, "Create a keystore", func(t *testing.T) {
		res, err := user.Client.CreateKeystoreWithResponse(s.Ctx,
			generated.CreateKeystoreJSONRequestBody{Name: keystoreName},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON201 != nil)
		assert.Equal(t, res.JSON201.Name, keystoreName)

		keystoreId = res.JSON201.Id
	})

	website, username, password, notes :=
		s.Random(s.T), s.Random(s.T), s.Random(s.T), s.Random(s.T)

	s.Run(s.T, "Add an entry to the keystore", func(t *testing.T) {
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
	})

	return keystoreId
}

func CreateInvitation(s *suite.Suite,
	keystoreId string, inviter *suite.Client, invitee *suite.Client,
) (invitationId string) {
	s.Run(s.T, "Invite someone to access the keystore", func(t *testing.T) {
		res, err := inviter.Client.CreateInvitationWithResponse(s.Ctx, keystoreId,
			generated.CreateInvitationJSONRequestBody{
				Invitee: invitee.Username,
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON201 != nil)

		invitationId = res.JSON201.Id
	})

	return invitationId
}
