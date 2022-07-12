package test

import (
	"context"
	"time"

	"myst/src/client/rest/generated"
)

func (s *IntegrationTestSuite) TestInvitations() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	keystoreName := s.rand()
	var keystoreId string

	s.Run("There is a keystore with one entry", func() {
		res, err := s.client1.CreateKeystoreWithResponse(ctx,
			generated.CreateKeystoreJSONRequestBody{Name: keystoreName},
		)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON201)
		s.Require().Equal(keystoreName, res.JSON201.Name)

		keystoreId = res.JSON201.Id
	})

	website := s.rand()
	username := s.rand()
	password := s.rand()
	notes := s.rand()

	s.Run("Add an entry to the keystore", func() {
		res, err := s.client1.CreateEntryWithResponse(ctx, keystoreId,
			generated.CreateEntryJSONRequestBody{
				Website:  website,
				Username: username,
				Password: password,
				Notes:    notes,
			},
		)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON201)

		s.Require().Equal(res.JSON201.Website, website)
		s.Require().Equal(res.JSON201.Username, username)
		s.Require().Equal(res.JSON201.Password, password)
		s.Require().Equal(res.JSON201.Notes, notes)
	})

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

	s.Run("The inviter will have access to the invitation", func() {
		res, err := s.client1.GetInvitationsWithResponse(ctx)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON200)
		s.Require().Len(*res.JSON200, 1)
		s.Require().Equal(invitationId, (*res.JSON200)[0].Id)
	})

}
