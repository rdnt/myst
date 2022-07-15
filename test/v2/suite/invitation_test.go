package suite

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/src/client/rest/generated"
)

func TestInvitations(t *testing.T) {
	s := setup(t)

	inviter := s.Client1
	invitee := s.Client2
	other := s.Client3

	keystoreId := s.CreateTestKeystore(t)

	var invitationId string
	s.Run(t, "Invite someone to access the keystore", func(ddt *testing.T) {
		res, err := inviter.client.CreateInvitationWithResponse(s.ctx,
			keystoreId,
			generated.CreateInvitationJSONRequestBody{
				Invitee: s.Client2.username,
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON201 != nil)

		invitationId = res.JSON201.Id
	})

	s.Run(t, "The inviter has access to the invitation", func(t *testing.T) {
		res, err := inviter.client.GetInvitationsWithResponse(s.ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 1)
		assert.Equal(t, (*res.JSON200)[0].Id, invitationId)

		res2, err := inviter.client.GetInvitationWithResponse(s.ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Id, invitationId)
		assert.Equal(t, res2.JSON200.Status, generated.Pending)
	})

	s.Run(t, "The invitee has access to the invitation", func(t *testing.T) {
		res, err := invitee.client.GetInvitationsWithResponse(s.ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 1)
		assert.Equal(t, (*res.JSON200)[0].Id, invitationId)

		res2, err := invitee.client.GetInvitationWithResponse(s.ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Id, invitationId)
		assert.Equal(t, res2.JSON200.Status, generated.Pending)
	})

	s.Run(t, "Another user doesn't have access to the invitation", func(t *testing.T) {
		res, err := other.client.GetInvitationsWithResponse(s.ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 0)

		res2, err := other.client.GetInvitationWithResponse(s.ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 == nil)
		assert.Assert(t, res2.JSONDefault != nil)
	})
}
