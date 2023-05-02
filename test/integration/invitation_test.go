package integration_test

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"

	"myst/src/client/rest/generated"
	"myst/test/integration/suite"
)

func TestInvitations(t *testing.T) {
	s := suite.New(t)

	err := s.Client1.App.Initialize(s.Client1.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client1.App.Register(s.Client1.Username, s.Client1.Password)
	assert.NilError(t, err)

	err = s.Client2.App.Initialize(s.Client2.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client2.App.Register(s.Client2.Username, s.Client2.Password)
	assert.NilError(t, err)

	err = s.Client3.App.Initialize(s.Client3.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client3.App.Register(s.Client3.Username, s.Client3.Password)
	assert.NilError(t, err)

	keystore := s.CreateKeystore(t)

	var invitationId string
	t.Run("Invitation is created", func(t *testing.T) {
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

	t.Run("The inviter has access to the invitation", func(t *testing.T) {
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

	t.Run("The invitee has access to the invitation", func(t *testing.T) {
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

	t.Run("Another user doesn't have access to the invitation", func(t *testing.T) {
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
	s := suite.New(t)

	err := s.Client1.App.Initialize(s.Client1.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client1.App.Register(s.Client1.Username, s.Client1.Password)
	assert.NilError(t, err)

	err = s.Client2.App.Initialize(s.Client2.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client2.App.Register(s.Client2.Username, s.Client2.Password)
	assert.NilError(t, err)

	keystore := s.CreateKeystore(t)
	invitationId := s.CreateInvitation(t, keystore.Id)

	t.Run("The inviter deletes the invitation", func(t *testing.T) {
		res, err := s.Client1.Client.DeclineOrCancelInvitationWithResponse(
			s.Ctx, invitationId,
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, res.JSON200.Status, generated.Deleted)
		assert.Assert(t, res.JSON200.DeletedAt != time.Time{})
	})

	t.Run("Only the inviter has access to the deleted invitation", func(t *testing.T) {
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
	s := suite.New(t)

	err := s.Client1.App.Initialize(s.Client1.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client1.App.Register(s.Client1.Username, s.Client1.Password)
	assert.NilError(t, err)

	err = s.Client2.App.Initialize(s.Client2.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client2.App.Register(s.Client2.Username, s.Client2.Password)
	assert.NilError(t, err)

	keystore := s.CreateKeystore(t)
	invitationId := s.CreateInvitation(t, keystore.Id)

	t.Run("The invitee declines the invitation", func(t *testing.T) {
		res, err := s.Client2.Client.DeclineOrCancelInvitationWithResponse(
			s.Ctx, invitationId,
		)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
		assert.Equal(t, res.JSON200.Status, generated.Declined)
		assert.Assert(t, res.JSON200.DeclinedAt != time.Time{})
	})

	t.Run("Both users see invitation as declined", func(t *testing.T) {
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
	s := suite.New(t)

	err := s.Client1.App.Initialize(s.Client1.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client1.App.Register(s.Client1.Username, s.Client1.Password)
	assert.NilError(t, err)

	err = s.Client2.App.Initialize(s.Client2.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client2.App.Register(s.Client2.Username, s.Client2.Password)
	assert.NilError(t, err)

	keystore := s.CreateKeystore(t)
	invitationId := s.CreateInvitation(t, keystore.Id)

	t.Run("The inviter cannot accept the invitation", func(t *testing.T) {
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

	t.Run("The invitee accepts the invitation", func(t *testing.T) {
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

	t.Run("Invitation cannot now be accepted by either user", func(t *testing.T) {
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
	s := suite.New(t)

	err := s.Client1.App.Initialize(s.Client1.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client1.App.Register(s.Client1.Username, s.Client1.Password)
	assert.NilError(t, err)

	err = s.Client2.App.Initialize(s.Client2.MasterPassword)
	assert.NilError(t, err)

	_, err = s.Client2.App.Register(s.Client2.Username, s.Client2.Password)
	assert.NilError(t, err)

	keystore := s.CreateKeystore(t)
	invitationId := s.CreateInvitation(t, keystore.Id)

	t.Run("The invitee accepts the invitation", func(t *testing.T) {
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

	t.Run("The inviter eventually finalizes the invitation", func(t *testing.T) {
		res, err := s.Client1.Client.FinalizeInvitationWithResponse(s.Ctx, invitationId)
		assert.NilError(t, err)
		assert.Assert(t, res.JSON200 != nil)
	})

	t.Run("Both users see invitation as finalized", func(t *testing.T) {
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
