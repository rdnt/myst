package suite

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/src/client/rest/generated"
	"myst/test/integration/suite"
)

func TestInvitations(t *testing.T) {
	s := suite.New(t)

	inviter := s.Client1
	invitee := s.Client2
	other := s.Client3

	keystoreId := s.CreateTestKeystore(t)

	var invitationId string
	s.Run(t, "Invite someone to access the keystore", func(ddt *testing.T) {
		res, err := inviter.Client.CreateInvitationWithResponse(s.Ctx,
			keystoreId,
			generated.CreateInvitationJSONRequestBody{
				Invitee: s.Client2.Username,
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON201 != nil)

		invitationId = res.JSON201.Id
	})

	s.Run(t, "The inviter has access to the invitation", func(t *testing.T) {
		res, err := inviter.Client.GetInvitationsWithResponse(s.Ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 1)
		assert.Equal(t, (*res.JSON200)[0].Id, invitationId)

		res2, err := inviter.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Id, invitationId)
		assert.Equal(t, res2.JSON200.Status, generated.Pending)
	})

	s.Run(t, "The invitee has access to the invitation", func(t *testing.T) {
		res, err := invitee.Client.GetInvitationsWithResponse(s.Ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 1)
		assert.Equal(t, (*res.JSON200)[0].Id, invitationId)

		res2, err := invitee.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Id, invitationId)
		assert.Equal(t, res2.JSON200.Status, generated.Pending)
	})

	s.Run(t, "Another user doesn't have access to the invitation", func(t *testing.T) {
		res, err := other.Client.GetInvitationsWithResponse(s.Ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 0)

		res2, err := other.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 == nil)
		assert.Assert(t, res2.JSONDefault != nil)
	})
}
