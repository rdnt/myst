package integration_test

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"

	"myst/src/client/rest/generated"
)

func TestInvitations(t *testing.T) {
	s := setup(t)

	keystore := s.createKeystore(t)

	var invitationId string
	s.Run(t, "Invitation is created", func(t *testing.T) {
		res, err := s.Client1.Client.CreateInvitationWithResponse(s.Ctx,
			keystore.Id,
			generated.CreateInvitationJSONRequestBody{
				Invitee: s.Client2.Username,
			},
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON201 != nil)

		invitationId = res.JSON201.Id
	})

	s.Run(t, "The inviter has access to the invitation", func(t *testing.T) {
		res, err := s.Client1.Client.GetInvitationsWithResponse(s.Ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 1)
		assert.Equal(t, (*res.JSON200)[0].Id, invitationId)

		res2, err := s.Client1.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Id, invitationId)
		assert.Equal(t, res2.JSON200.Status, generated.Pending)
	})

	s.Run(t, "The invitee has access to the invitation", func(t *testing.T) {
		res, err := s.Client2.Client.GetInvitationsWithResponse(s.Ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 1)
		assert.Equal(t, (*res.JSON200)[0].Id, invitationId)

		res2, err := s.Client2.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Id, invitationId)
		assert.Equal(t, res2.JSON200.Status, generated.Pending)
	})

	s.Run(t, "Another user doesn't have access to the invitation", func(t *testing.T) {
		res, err := s.Client3.Client.GetInvitationsWithResponse(s.Ctx)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, len(*res.JSON200) == 0)

		res2, err := s.Client3.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 == nil)
		assert.Assert(t, res2.JSONDefault != nil)
	})
}

func TestInvitationDelete(t *testing.T) {
	s := setup(t)

	keystore := s.createKeystore(t)
	invitationId := s.createInvitation(t, keystore.Id)

	s.Run(t, "The inviter deletes the invitation", func(t *testing.T) {
		res, err := s.Client1.Client.DeclineOrCancelInvitationWithResponse(
			s.Ctx, invitationId,
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, res.JSON200.Status, generated.Deleted)
		assert.Assert(t, res.JSON200.DeletedAt != time.Time{})
	})

	s.Run(t, "Only the inviter has access to the deleted invitation", func(t *testing.T) {
		res, err := s.Client1.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, res.JSON200.Status, generated.Deleted)
		assert.Assert(t, res.JSON200.DeletedAt != time.Time{})

		res, err = s.Client2.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 == nil)
		assert.Assert(t, res.JSONDefault != nil)
	})
}

func TestInvitationDecline(t *testing.T) {
	s := setup(t)

	keystore := s.createKeystore(t)
	invitationId := s.createInvitation(t, keystore.Id)

	s.Run(t, "The invitee declines the invitation", func(t *testing.T) {
		res, err := s.Client2.Client.DeclineOrCancelInvitationWithResponse(
			s.Ctx, invitationId,
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, res.JSON200.Status, generated.Declined)
		assert.Assert(t, res.JSON200.DeclinedAt != time.Time{})
	})

	s.Run(t, "Both users see invitation as declined", func(t *testing.T) {
		res, err := s.Client1.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, res.JSON200.Status, generated.Declined)

		res, err = s.Client2.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, res.JSON200.Status, generated.Declined)
	})
}
func TestInvitationAccept(t *testing.T) {
	s := setup(t)

	keystore := s.createKeystore(t)
	invitationId := s.createInvitation(t, keystore.Id)

	s.Run(t, "The inviter cannot accept the invitation", func(t *testing.T) {
		res, err := s.Client1.Client.AcceptInvitationWithResponse(
			s.Ctx, invitationId,
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 == nil)
		assert.Assert(t, res.JSONDefault != nil)

		res2, err := s.Client1.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Status, generated.Pending)
	})

	s.Run(t, "The invitee accepts the invitation", func(t *testing.T) {
		res, err := s.Client2.Client.AcceptInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)

		res2, err := s.Client1.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Status, generated.Accepted)

		res2, err = s.Client2.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Status, generated.Accepted)
	})

	s.Run(t, "Invitation cannot now be accepted by either user", func(t *testing.T) {
		res, err := s.Client1.Client.AcceptInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSONDefault != nil)
		assert.Assert(t, res.JSON200 == nil)

		res2, err := s.Client1.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Status, generated.Accepted)

		res, err = s.Client2.Client.AcceptInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSONDefault != nil)
		assert.Assert(t, res.JSON200 == nil)

		res2, err = s.Client1.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Status, generated.Accepted)
	})
}

func TestInvitationFinalize(t *testing.T) {
	s := setup(t)

	keystore := s.createKeystore(t)
	invitationId := s.createInvitation(t, keystore.Id)

	s.Run(t, "The invitee accepts the invitation", func(t *testing.T) {
		res, err := s.Client2.Client.AcceptInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)

		res2, err := s.Client1.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Status, generated.Accepted)

		res2, err = s.Client2.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res2.JSON200 != nil)
		assert.Equal(t, res2.JSON200.Status, generated.Accepted)
	})

	s.Run(t, "The inviter eventually finalizes the invitation", func(t *testing.T) {
		res, err := s.Client1.Client.FinalizeInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
	})

	s.Run(t, "Both users see invitation as finalized", func(t *testing.T) {
		res, err := s.Client1.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, res.JSON200.Status == generated.Finalized)

		res, err = s.Client2.Client.GetInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Assert(t, res.JSON200.Status == generated.Finalized)
	})
}
