package client

import (
	"context"
	"fmt"

	"myst/internal/server/api/http/generated"

	"github.com/pkg/errors"
)

func (c *client) Invitation(keystoreId, invitationId string) (*generated.Invitation, error) {
	res, err := c.client.InvitationWithResponse(context.Background(), invitationId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get invitation")
	}

	if res.JSON200 == nil {
		return nil, fmt.Errorf("invalid response")
	}

	return res.JSON200, nil
}

//func (c *client) Invitations() ([]generated.Invitation, error) {
//	res, err := c.client.InvitationsWithResponse(context.Background())
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to get invitations")
//	}
//
//	if res.JSON200 == nil {
//		return nil, fmt.Errorf("invalid response")
//	}
//
//	return *res.JSON200, nil
//}

//func (c *client) CreateInvitation(keystoreId, inviteeId string, inviterPublicKey []byte) (*generated.Invitation, error) {
//	if !c.SignedIn() {
//		return nil, ErrNotSignedIn
//	}
//
//	res, err := c.client.CreateInvitationWithResponse(
//		context.Background(), keystoreId, generated.CreateInvitationJSONRequestBody{
//			InviteeId: inviteeId,
//			PublicKey: inviterPublicKey,
//		},
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	if res.JSON200 == nil {
//		return nil, fmt.Errorf("invalid response")
//	}
//
//	return res.JSON200, nil
//}

func (c *client) AcceptInvitation(keystoreId, invitationId string, inviteePublicKey []byte) (*generated.Invitation, error) {
	if !c.SignedIn() {
		return nil, ErrNotSignedIn
	}

	res, err := c.client.AcceptInvitationWithResponse(
		context.Background(), keystoreId, invitationId, generated.AcceptInvitationJSONRequestBody{
			PublicKey: inviteePublicKey,
		},
	)
	if err != nil {
		return nil, err
	}

	if res.JSON200 == nil {
		return nil, fmt.Errorf("invalid response")
	}

	return res.JSON200, nil
}

func (c *client) FinalizeInvitation(keystoreId, invitationId string, encryptedKeystoreKey []byte) (*generated.Invitation, error) {
	if !c.SignedIn() {
		return nil, ErrNotSignedIn
	}

	res2, err := c.client.FinalizeInvitationWithResponse(
		context.Background(), keystoreId, invitationId, generated.FinalizeInvitationJSONRequestBody{
			KeystoreKey: encryptedKeystoreKey,
		},
	)
	if err != nil {
		return nil, err
	}

	if res2.JSON200 == nil {
		return nil, fmt.Errorf("invalid response")
	}

	return res2.JSON200, nil
}
