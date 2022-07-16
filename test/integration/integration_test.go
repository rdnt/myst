package integration_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/src/client/rest/generated"
	"myst/test/integration/suite"
)

type Suite struct {
	*suite.Suite
}

func setup(t *testing.T) *Suite {
	return &Suite{suite.New(t)}
}

func (s *Suite) createKeystore(t *testing.T) (keystoreId string) {
	keystoreName := s.Random(t)

	s.Run(t, "Create a keystore", func(t *testing.T) {
		res, err := s.Client1.Client.CreateKeystoreWithResponse(s.Ctx,
			generated.CreateKeystoreJSONRequestBody{Name: keystoreName},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON201 != nil)
		assert.Equal(t, res.JSON201.Name, keystoreName)

		keystoreId = res.JSON201.Id
	})

	website, username, password, notes :=
		s.Random(t), s.Random(t), s.Random(t), s.Random(t)

	s.Run(t, "Add an entry to the keystore", func(t *testing.T) {
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
	})

	return keystoreId
}

func (s *Suite) createInvitation(t *testing.T, keystoreId string) (invitationId string) {
	s.Run(t, "Invite someone to access the keystore", func(t *testing.T) {
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
