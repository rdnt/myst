package remote

import (
	"context"
	"fmt"
	"myst/internal/server/api/http/generated"
	"myst/pkg/crypto"

	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"
)

func (r *remote) CreateInvitation(keystoreId, inviteeId string) (*generated.Invitation, error) {
	if r.bearerToken == "" {
		return nil, fmt.Errorf("not signed in")
	}

	res, err := r.client.CreateInvitationWithResponse(
		context.Background(), keystoreId, generated.CreateInvitationJSONRequestBody{
			InviteeId: inviteeId,
			PublicKey: r.publicKey,
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

func (r *remote) AcceptInvitation(keystoreId, invitationId string) (*generated.Invitation, error) {
	if r.bearerToken == "" {
		return nil, fmt.Errorf("not signed in")
	}

	//res, err := r.client.InvitationWithResponse(context.Background(), keystoreId, invitationId)
	//if err != nil {
	//	return nil, errors.Wrap(err, "failed to get invitation")
	//}
	//
	//if res.JSON200 == nil || res.JSON200.InviterKey == nil {
	//	return nil, fmt.Errorf("invalid response")
	//}
	//
	//// TODO: not needed for acceptance, needed after finalization to decrypt the keystore key
	//inviterKey := *res.JSON200.InviterKey
	//
	//asymKey, err := curve25519.X25519(r.privateKey, inviterKey)
	//if err != nil {
	//	panic(err)
	//}

	res, err := r.client.AcceptInvitationWithResponse(
		context.Background(), keystoreId, invitationId, generated.AcceptInvitationJSONRequestBody{
			PublicKey: r.publicKey,
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

func (r *remote) FinalizeInvitation(keystoreId, invitationId string, keystoreKey []byte) (*generated.Invitation, error) {
	if r.bearerToken == "" {
		return nil, fmt.Errorf("not signed in")
	}

	res, err := r.client.InvitationWithResponse(context.Background(), keystoreId, invitationId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get invitation")
	}
	if res.JSON200 == nil || res.JSON200.InviteeKey == nil {
		return nil, fmt.Errorf("invalid response")
	}

	inviteeKey := *res.JSON200.InviteeKey

	asymKey, err := curve25519.X25519(r.privateKey, inviteeKey)
	if err != nil {
		panic(err)
	}

	// encrypt the keystore key with the asymmetric key
	encryptedKeystoreKey, err := crypto.AES256CBC_Encrypt(asymKey, keystoreKey)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to encrypt keystore key")
	}

	res2, err := r.client.FinalizeInvitationWithResponse(
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
