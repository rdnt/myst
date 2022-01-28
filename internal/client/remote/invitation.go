package remote

import (
	"context"
	"encoding/hex"
	"fmt"

	"myst/internal/client/application/domain/invitation"
	"myst/internal/server/api/http/generated"
)

func (r *remote) CreateInvitation(keystoreId, inviteeId string, publicKey []byte) (*invitation.Invitation, error) {
	fmt.Println("CreateInvitation", keystoreId, inviteeId, publicKey)

	if r.bearerToken == "" {
		return nil, fmt.Errorf("not signed in")
	}

	res, err := r.client.CreateInvitationWithResponse(
		context.Background(), keystoreId, generated.CreateInvitationJSONRequestBody{
			InviteeId: inviteeId,
			PublicKey: hex.EncodeToString(publicKey),
		},
		r.authenticate(),
	)
	if err != nil {
		return nil, err
	}

	return r.parseInvitation(res.JSON200)
}

func (r *remote) AcceptInvitation(keystoreId, invitationId string, publicKey []byte) (*invitation.Invitation, error) {
	fmt.Println("AcceptInvitation", invitationId, publicKey)

	if r.bearerToken == "" {
		return nil, fmt.Errorf("not signed in")
	}

	res, err := r.client.AcceptInvitationWithResponse(
		context.Background(), keystoreId, invitationId, generated.AcceptInvitationJSONRequestBody{
			PublicKey: hex.EncodeToString(publicKey),
		},
		r.authenticate(),
	)
	if err != nil {
		return nil, err
	}

	return r.parseInvitation(res.JSON200)
}

func (r *remote) FinalizeInvitation(keystoreId, invitationId string, keystoreKey []byte) (
	*invitation.Invitation, error,
) {
	fmt.Println("FinalizeInvitation", invitationId, keystoreKey)

	if r.bearerToken == "" {
		return nil, fmt.Errorf("not signed in")
	}

	res, err := r.client.FinalizeInvitationWithResponse(
		context.Background(), keystoreId, invitationId, generated.FinalizeInvitationJSONRequestBody{
			KeystoreKey: hex.EncodeToString(keystoreKey),
		},
		r.authenticate(),
	)
	if err != nil {
		return nil, err
	}

	return r.parseInvitation(res.JSON200)
}

func (r *remote) parseInvitation(gen *generated.Invitation) (*invitation.Invitation, error) {
	if gen == nil {
		return nil, ErrInvalidResponse
	}

	return invitation.New(
		invitation.WithId(gen.Id),
	)
}
