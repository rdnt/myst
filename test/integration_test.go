package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	// servergen "myst/internal/server/api/rest/generated"
	"myst/pkg/testing/capture"
	"myst/src/client/rest/generated"
)

func TestIntegration(t *testing.T) {
	s := &IntegrationTestSuite{
		capture: capture.New(t),
	}
	suite.Run(t, s)
}

func (s *IntegrationTestSuite) createTestKeystore(ctx context.Context) (keystoreId string) {
	keystoreName := s.rand()

	s.Run("Create a keystore", func() {
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

	return keystoreId
}

func (s *IntegrationTestSuite) createDebugInvitation(ctx context.Context,
	keystoreId string, inviter, invitee string,
) (invitationId string) {
	s.Run("Invite someone to access the keystore", func() {
		res, err := inviter.CreateInvitationWithResponse(ctx, keystoreId,
			generated.CreateInvitationJSONRequestBody{
				Invitee: invitee,
			},
		)
		s.Require().Nil(err)
		s.Require().NotNil(res.JSON201)

		invitationId = res.JSON201.Id
	})

	return invitationId
}
