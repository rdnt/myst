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

	keystoreId := CreateKeystore(s)

	var invitationId string
	s.Run("Invite someone to access the keystore", func(ddt *testing.T) {
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

	s.Run("The inviter has access to the invitation", func(t *testing.T) {
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

	s.Run("The invitee has access to the invitation", func(t *testing.T) {
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

	s.Run("Another user doesn't have access to the invitation", func(t *testing.T) {
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

func TestInvitationAccept(t *testing.T) {
	s := suite.New(t)

	inviter := s.Client1
	invitee := s.Client2
	other := s.Client3

	keystoreId := CreateKeystore(s)

	invitationId := CreateInvitation(s, keystoreId, inviter, invitee)

	s.Run("The inviter cannot accept the invitation", func(t *testing.T) {
		res, err := inviter.Client.AcceptInvitationWithResponse(
			s.Ctx, invitationId,
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 == nil)
		assert.Assert(t, res.JSONDefault != nil)

		res2, err := inviter.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Status, generated.Pending)
	})

	s.Run("The invitee accepts the invitation", func(t *testing.T) {
		res, err := invitee.Client.AcceptInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)

		res2, err := inviter.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Status, generated.Accepted)

		res2, err = invitee.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Status, generated.Accepted)
	})

	s.Run("Invitation cannot be accepted by either user", func(t *testing.T) {
		res, err := inviter.Client.AcceptInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSONDefault != nil)
		assert.Assert(t, res.JSON200 == nil)

		res2, err := inviter.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Status, generated.Accepted)

		res, err = invitee.Client.AcceptInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSONDefault != nil)
		assert.Assert(t, res.JSON200 == nil)

		res2, err = inviter.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Status, generated.Accepted)
	})

	// s.Run("The inviter eventually finalizes the invitation", func() {
	// 	// manually trigger finalization
	// 	// TODO: maybe do eventually with retries
	// 	_, err := s._client1.app.FinalizeInvitation(invitationId)
	// 	s.Require().Nil(err)
	// })
	//
	// s.Run("Invitation is finalized", func() {
	// 	res, err := s.client1.GetInvitationWithResponse(ctx, invitationId)
	// 	s.Require().Nil(err)
	// 	s.Require().NotNil(res.JSON200)
	// 	s.Require().Equal(res.JSON200.Status, generated.Finalized)
	//
	// 	res, err = s.client2.GetInvitationWithResponse(ctx, invitationId)
	// 	s.Require().Nil(err)
	// 	s.Require().NotNil(res.JSON200)
	// 	s.Require().Equal(res.JSON200.Status, generated.Finalized)
	// })
}
