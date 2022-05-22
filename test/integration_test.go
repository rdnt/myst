package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	clientgen "myst/internal/client/api/http/generated"
	"myst/pkg/optional"

	//servergen "myst/internal/server/api/http/generated"
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

	masterPassword := "12345678"

	createksres, err := s.client1.CreateKeystoreWithResponse(ctx, clientgen.CreateKeystoreJSONRequestBody{
		Name:     "test-keystore-1",
		Password: optional.Ref(masterPassword),
	})
	s.Require().NoError(err)
	s.Require().NotNil(createksres.JSON201)

	ks := *createksres.JSON201
	s.Require().Equal(ks.Name, "test-keystore-1")

	createinvres, err := s.client1.CreateInvitationWithResponse(ctx, ks.Id, clientgen.CreateInvitationJSONRequestBody{
		InviteeId: "abcd",
	})
	s.Require().NoError(err)
	s.Require().NotNil(createinvres.JSON200)

	inv := *createinvres.JSON200
	s.T().Log(inv)

	invres, err := s.client1.GetInvitationsWithResponse(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(invres.JSON200)
	s.Require().Len(*invres.JSON200, 1)

	invres, err = s.client2.GetInvitationsWithResponse(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(invres.JSON200)
	s.Require().Len(*invres.JSON200, 1)

	acceptResponse, err := s.client2.AcceptInvitationWithResponse(ctx, inv.Id)
	s.Require().NoError(err)
	s.Require().NotNil(acceptResponse.JSON200)

	inv = *acceptResponse.JSON200
	s.Require().Equal(inv.Status, clientgen.InvitationStatusAccepted)
	s.Require().NotNil(inv.InviteeKey)

	finalizedInv, err := s._client1.app.FinalizeInvitation(inv.Id)
	s.Require().NoError(err)
	s.Require().NotNil(acceptResponse.JSON200)
	s.Require().Equal(finalizedInv.Id, inv.Id)

	s.T().Log(finalizedInv)

	// TODO: accept/decline keystore invitation
}
