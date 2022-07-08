package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	clientgen "myst/internal/client/api/rest/generated"
	"myst/pkg/optional"

	//servergen "myst/internal/server/api/rest/generated"
	"myst/pkg/testing/capture"
)

func TestIntegration(t *testing.T) {
	t.Log("Integration test")
	s := &IntegrationTestSuite{
		capture: capture.New(t),
	}
	suite.Run(t, s)
}

func (s *IntegrationTestSuite) TestKeystoreCreation() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	masterPassword := "12345678"
	k1name := "my-keystore-1"

	createksres, err := s.client1.CreateKeystoreWithResponse(ctx, clientgen.CreateKeystoreJSONRequestBody{
		Name:     k1name,
		Password: optional.Ref(masterPassword),
	})
	s.Require().NoError(err)
	s.Require().NotNil(createksres.JSON201)

	ks := *createksres.JSON201
	s.Require().Equal(ks.Name, k1name)
	s.T().Logf("Keystore created: %#v\n", ks)

	_, err = s.client1.CreateEntryWithResponse(ctx, ks.Id, clientgen.CreateEntryJSONRequestBody{
		Website:  "example.com",
		Username: "xmpl",
		Password: "1337",
		Notes:    "blah blah blah",
	})
	s.Require().NoError(err)

	restKeystore, err := s.client1.KeystoreWithResponse(ctx, ks.Id)
	s.Require().NoError(err)
	s.Require().NotNil(restKeystore.JSON200)

	s.T().Log("@@@@", *restKeystore.JSON200)

	masterPassword2 := "87654321"
	k2name := "my-keystore-2"

	createksres2, err := s.client2.CreateKeystoreWithResponse(ctx, clientgen.CreateKeystoreJSONRequestBody{
		Name:     k2name,
		Password: optional.Ref(masterPassword2),
	})
	s.Require().NoError(err)
	s.Require().NotNil(createksres2.JSON201)

	ks2 := *createksres2.JSON201
	s.Require().Equal(ks2.Name, k2name)

	// ###

	res, err := s.server.KeystoresWithResponse(context.Background())
	s.Require().NoError(err)
	fmt.Println("@@@@@@@!!!!!!!!!!!!!!!!!!!!", string(res.Body))

	// ###

	s.T().Log("Creating invitation", ks.Id)
	createinvres, err := s.client1.CreateInvitationWithResponse(ctx, ks.Id, clientgen.CreateInvitationJSONRequestBody{
		InviteeId: "abcd",
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
	s.Require().Equal(restInv.Status, clientgen.InvitationStatusAccepted)
	s.Require().NotNil(restInv.InviteeKey)

	fmt.Println("INVID", restInv.Id)

	inv, err := s._client1.app.FinalizeInvitation(restInv.Id)
	s.Require().NoError(err)
	s.Require().NotNil(acceptResponse.JSON200)
	s.Require().Equal(inv.Id, restInv.Id)

	s.T().Log(inv)

	// TODO: accept/decline keystore invitation

	//time.Sleep(10 * time.Minute)
}
