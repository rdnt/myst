package test

import (
	"context"
	"time"

	"myst/src/client/rest/generated"
)

func (s *IntegrationTestSuite) TestInvitations() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	keystoreId := s.createTestKeystore(ctx)

	var invitationId string

	s.Run("Invite someone to access the keystore", func() {
		res, err := s.client1.CreateInvitationWithResponse(ctx, keystoreId,
			generated.CreateInvitationJSONRequestBody{
				Invitee: s._client2.username,
			},
		)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON201)

		invitationId = res.JSON201.Id
	})

	s.Run("The inviter has access to the invitation", func() {
		res, err := s.client1.GetInvitationsWithResponse(ctx)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)
		s.Require().Len(*res.JSON200, 1)
		s.Require().Equal(invitationId, (*res.JSON200)[0].Id)

		res2, err := s.client1.GetInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res2.JSON200)
		s.Require().Equal(res2.JSON200.Id, (*res.JSON200)[0].Id)
		s.Require().Equal(res2.JSON200.Status, generated.Pending)
	})

	s.Run("The invitee has access to the invitation", func() {
		res, err := s.client2.GetInvitationsWithResponse(ctx)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)
		s.Require().Len(*res.JSON200, 1)
		s.Require().Equal(invitationId, (*res.JSON200)[0].Id)

		res2, err := s.client2.GetInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res2.JSON200)
		s.Require().Equal(res2.JSON200.Id, (*res.JSON200)[0].Id)
		s.Require().Equal(res2.JSON200.Status, generated.Pending)
	})

	s.Run("Another user doesn't have access to the invitation", func() {
		res, err := s.client3.GetInvitationsWithResponse(ctx)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)
		s.Require().Len(*res.JSON200, 0)

		res2, err := s.client3.GetInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res2.JSONDefault)
	})
}

func (s *IntegrationTestSuite) TestInvitationAccept() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	keystoreId := s.createTestKeystore(ctx)
	invitationId := s.createDebugInvitation(ctx, keystoreId, s._client2.username)

	s.Run("The inviter cannot accept the invitation", func() {
		res, err := s.client1.AcceptInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSONDefault)

		res2, err := s.client1.GetInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res2.JSON200)
		s.Require().Equal(res2.JSON200.Status, generated.Pending)
	})

	s.Run("The invitee accepts the invitation", func() {
		res, err := s.client2.AcceptInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)

		res2, err := s.client1.GetInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res2.JSON200)
		s.Require().Equal(res2.JSON200.Status, generated.Accepted)

		res2, err = s.client2.GetInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res2.JSON200)
		s.Require().Equal(res2.JSON200.Status, generated.Accepted)
	})

	s.Run("Invitation cannot be accepted by either user at this point", func() {
		res, err := s.client1.AcceptInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSONDefault)

		res2, err := s.client1.GetInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res2.JSON200)
		s.Require().Equal(res2.JSON200.Status, generated.Accepted)

		res, err = s.client2.AcceptInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSONDefault)

		res2, err = s.client1.GetInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res2.JSON200)
		s.Require().Equal(res2.JSON200.Status, generated.Accepted)
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

func (s *IntegrationTestSuite) TestInvitationDelete() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	keystoreId := s.createTestKeystore(ctx)
	invitationId := s.createDebugInvitation(ctx, keystoreId, s._client2.username)

	s.Run("Both users can see the invitation", func() {
		res, err := s.client1.GetInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)
		s.Require().Equal(res.JSON200.Status, generated.Pending)

		res, err = s.client2.GetInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)
		s.Require().Equal(res.JSON200.Status, generated.Pending)
	})

	s.Run("The inviter deletes the invitation", func() {
		res, err := s.client1.DeclineOrCancelInvitationWithResponse(
			ctx, invitationId,
		)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)
	})

	s.Run("Only the inviter has access to the deleted invitation", func() {
		res, err := s.client1.GetInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)

		res, err = s.client2.GetInvitationWithResponse(ctx, invitationId)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSONDefault)
	})
	//
	// s.Run("The invitee accepts the invitation", func() {
	// 	res, err := s.client2.AcceptInvitationWithResponse(ctx, invitationId)
	// 	s.Require().Nil(err)
	// 	s.Require().NotNil(res.JSON200)
	//
	// 	res2, err := s.client1.GetInvitationWithResponse(ctx, invitationId)
	// 	s.Require().Nil(err)
	// 	s.Require().NotNil(res2.JSON200)
	// 	s.Require().Equal(res2.JSON200.Status, generated.Accepted)
	//
	// 	res2, err = s.client2.GetInvitationWithResponse(ctx, invitationId)
	// 	s.Require().Nil(err)
	// 	s.Require().NotNil(res2.JSON200)
	// 	s.Require().Equal(res2.JSON200.Status, generated.Accepted)
	// })
	//
	// s.Run("Invitation cannot be accepted by either user at this point", func() {
	// 	res, err := s.client1.AcceptInvitationWithResponse(ctx, invitationId)
	// 	s.Require().Nil(err)
	// 	s.Require().NotNil(res.JSONDefault)
	//
	// 	res2, err := s.client1.GetInvitationWithResponse(ctx, invitationId)
	// 	s.Require().Nil(err)
	// 	s.Require().NotNil(res2.JSON200)
	// 	s.Require().Equal(res2.JSON200.Status, generated.Accepted)
	//
	// 	res, err = s.client2.AcceptInvitationWithResponse(ctx, invitationId)
	// 	s.Require().Nil(err)
	// 	s.Require().NotNil(res.JSONDefault)
	//
	// 	res2, err = s.client1.GetInvitationWithResponse(ctx, invitationId)
	// 	s.Require().Nil(err)
	// 	s.Require().NotNil(res2.JSON200)
	// 	s.Require().Equal(res2.JSON200.Status, generated.Accepted)
	// })

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
