package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"myst/pkg/optional"
	clientgen "myst/src/client/rest/generated"

	// servergen "myst/internal/server/api/rest/generated"
	"myst/pkg/testing/capture"
)

func TestIntegration(t *testing.T) {
	s := &IntegrationTestSuite{
		capture: capture.New(t),
	}
	suite.Run(t, s)
}

func (s *IntegrationTestSuite) TestKeystoreCreation() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	k1name := s.rand()

	createksres, err := s.client1.CreateKeystoreWithResponse(ctx, clientgen.CreateKeystoreJSONRequestBody{
		Name:     k1name,
		Password: optional.Ref(s._client1.masterPassword),
	})
	s.Require().NoError(err)
	s.Require().NotNil(createksres.JSON201)

	ks := *createksres.JSON201
	s.Require().Equal(ks.Name, k1name)

	_, err = s.client1.CreateEntryWithResponse(ctx, ks.Id, clientgen.CreateEntryJSONRequestBody{
		Website:  s.rand(),
		Username: s.rand(),
		Password: s.rand(),
		Notes:    s.rand(),
	})
	s.Require().NoError(err)

	restKeystore, err := s.client1.KeystoreWithResponse(ctx, ks.Id)
	s.Require().NoError(err)
	s.Require().NotNil(restKeystore.JSON200)

	k2name := s.rand()

	createksres2, err := s.client2.CreateKeystoreWithResponse(ctx, clientgen.CreateKeystoreJSONRequestBody{
		Name:     k2name,
		Password: optional.Ref(s._client2.masterPassword),
	})
	s.Require().NoError(err)
	s.Require().NotNil(createksres2.JSON201)

	ks2 := *createksres2.JSON201
	s.Require().Equal(ks2.Name, k2name)

	_, err = s.server.KeystoresWithResponse(context.Background())
	s.Require().NoError(err)

	createinvres, err := s.client1.CreateInvitationWithResponse(ctx, ks.Id, clientgen.CreateInvitationJSONRequestBody{
		Invitee: s._client2.username,
	})
	s.Require().NoError(err)
	s.Require().NotNil(createinvres.JSON200)

	restInv := *createinvres.JSON200

	invres, err := s.client1.GetInvitationsWithResponse(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(invres.JSON200)
	s.Require().Len(*invres.JSON200, 1)

	invres, err = s.client2.GetInvitationsWithResponse(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(invres.JSON200)
	s.Require().Len(*invres.JSON200, 1)

	acceptResponse, err := s.client2.AcceptInvitationWithResponse(ctx, restInv.Id)
	s.Require().NoError(err)
	s.Require().NotNil(acceptResponse.JSON200)

	restInv = *acceptResponse.JSON200
	s.Require().Equal(restInv.Status, clientgen.Accepted)
	s.Require().NotNil(restInv.InviteePublicKey)

	inv, err := s._client1.app.FinalizeInvitation(restInv.Id)
	s.Require().NoError(err)
	s.Require().NotNil(acceptResponse.JSON200)
	s.Require().Equal(inv.Id, restInv.Id)
}
