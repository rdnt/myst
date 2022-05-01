package remote

import (
	"context"
	"encoding/hex"
	"fmt"

	"myst/internal/server/api/http/generated"
)

func (r *remote) CreateInvitation(keystoreId, inviteeId string, publicKey []byte) (*generated.Invitation, error) {
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

	if res.JSON200 == nil {
		return nil, fmt.Errorf("invalid response")
	}

	return res.JSON200, nil
}

func (r *remote) AcceptInvitation(keystoreId, invitationId string, publicKey []byte) (*generated.Invitation, error) {
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

	if res.JSON200 == nil {
		return nil, fmt.Errorf("invalid response")
	}

	return res.JSON200, nil
}

func (r *remote) FinalizeInvitation(keystoreId, invitationId string, keystoreKey []byte) (*generated.Invitation, error) {
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

	if res.JSON200 == nil {
		return nil, fmt.Errorf("invalid response")
	}

	return res.JSON200, nil
}
