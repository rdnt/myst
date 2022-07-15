package suite

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"

	"myst/src/client/rest/generated"
	"myst/test/integration/suite"
)

func TestInvitations(t *testing.T) {
	s := suite.New(t)

	inviter := s.Client1
	invitee := s.Client2
	other := s.Client3

	keystoreId := CreateKeystore(t, s)

	var invitationId string
	s.Run("Invitation is created", func(ddt *testing.T) {
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

func TestInvitationDelete(t *testing.T) {
	s := suite.New(t)

	inviter := s.Client1
	invitee := s.Client2

	keystoreId := CreateKeystore(t, s)
	invitationId := CreateInvitation(s, keystoreId, inviter, invitee)

	s.Run("The inviter deletes the invitation", func(t *testing.T) {
		res, err := inviter.Client.DeclineOrCancelInvitationWithResponse(
			s.Ctx, invitationId,
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, res.JSON200.Status, generated.Deleted)
		assert.Assert(t, res.JSON200.DeletedAt != time.Time{})
	})

	s.Run("Only the inviter has access to the deleted invitation", func(t *testing.T) {
		res, err := inviter.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, res.JSON200.Status, generated.Deleted)
		assert.Assert(t, res.JSON200.DeletedAt != time.Time{})

		res, err = invitee.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 == nil)
		assert.Assert(t, res.JSONDefault != nil)
	})
}

func TestInvitationDecline(t *testing.T) {
	s := suite.New(t)

	inviter := s.Client1
	invitee := s.Client2

	keystoreId := CreateKeystore(t, s)
	invitationId := CreateInvitation(s, keystoreId, inviter, invitee)

	s.Run("The invitee declines the invitation", func(t *testing.T) {
		res, err := invitee.Client.DeclineOrCancelInvitationWithResponse(
			s.Ctx, invitationId,
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, res.JSON200.Status, generated.Declined)
		assert.Assert(t, res.JSON200.DeclinedAt != time.Time{})
	})

	s.Run("Both users see invitation as declined", func(t *testing.T) {
		res, err := inviter.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, res.JSON200.Status, generated.Declined)

		res, err = invitee.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, res.JSON200.Status, generated.Declined)
	})
}
func TestInvitationAccept(t *testing.T) {
	s := suite.New(t)

	inviter := s.Client1
	invitee := s.Client2

	keystoreId := CreateKeystore(t, s)
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

	s.Run("Invitation cannot now be accepted by either user", func(t *testing.T) {
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
}

func TestInvitationFinalize(t *testing.T) {
	s := suite.New(t)

	inviter := s.Client1
	invitee := s.Client2

	keystoreId := CreateKeystore(t, s)
	invitationId := CreateInvitation(s, keystoreId, inviter, invitee)

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

	s.Run("The inviter eventually finalizes the invitation", func(t *testing.T) {
		// manually trigger finalization
		// TODO: maybe do eventually with retries
		_, err := inviter.App.FinalizeInvitation(invitationId)
		assert.NilError(t, err)
	})

	s.Run("Both users see invitation as finalized", func(t *testing.T) {
		res, err := inviter.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, res.JSON200.Status == generated.Finalized)

		res, err = invitee.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, res.JSON200.Status == generated.Finalized)
	})
}
