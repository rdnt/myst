package remote

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"myst/internal/client/application/domain/invitation"
	"myst/internal/server/api/http/generated"
)

var (
	ErrSignedOut = errors.New("signed out")
)

func (r *remote) Invitation(id string) (invitation.Invitation, error) {
	res, err := r.client.InvitationWithResponse(context.Background(), id)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to get invitation")
	}

	// TODO: add statuscode check

	if res.JSON200 == nil {
		return invitation.Invitation{}, fmt.Errorf("invalid response")
	}

	inv, err := InvitationFromJSON(*res.JSON200)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *remote) UpdateInvitation(k invitation.Invitation) error {
	panic("implement me")
}

func (r *remote) DeleteInvitation(id string) error {
	panic("implement me")
}

func (r *remote) CreateInvitation(inv invitation.Invitation) (invitation.Invitation, error) {
	if !r.SignedIn() {
		return invitation.Invitation{}, ErrSignedOut
	}

	res, err := r.client.CreateInvitationWithResponse(
		context.Background(), inv.KeystoreId, generated.CreateInvitationJSONRequestBody{
			InviteeId: inv.InviteeId,
			PublicKey: r.publicKey,
		},
	)
	if err != nil {
		return invitation.Invitation{}, err
	}

	if res.JSON200 == nil {
		return invitation.Invitation{}, fmt.Errorf("invalid response")
	}

	inv, err = InvitationFromJSON(*res.JSON200)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *remote) Invitations() (map[string]invitation.Invitation, error) {
	res, err := r.client.InvitationsWithResponse(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get invitations")
	}

	if res.JSON200 == nil {
		return nil, fmt.Errorf("invalid response")
	}

	invs := map[string]invitation.Invitation{}
	for _, inv := range *res.JSON200 {
		inv, err := InvitationFromJSON(inv)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to parse invitation")
		}

		invs[inv.Id] = inv
	}

	return invs, nil
}

//func (r *remote) AcceptInvitation(keystoreId, invitationId string) (invitation.Invitation, error) {
//	if !r.SignedIn() {
//		return invitation.Invitation{}, ErrSignedOut
//	}
//
//	res, err := r.client.AcceptInvitationWithResponse(
//		context.Background(), keystoreId, invitationId, generated.AcceptInvitationJSONRequestBody{
//			PublicKey: r.publicKey,
//		},
//	)
//	if err != nil {
//		return invitation.Invitation{}, err
//	}
//
//	if res.JSON200 == nil {
//		return invitation.Invitation{}, fmt.Errorf("invalid response")
//	}
//
//	inv, err := InvitationFromJSON(*res.JSON200)
//	if err != nil {
//		return invitation.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
//	}
//
//	return inv, nil
//}
//
//func (r *remote) FinalizeInvitation(localKeystoreId, keystoreId, invitationId string) (invitation.Invitation, error) {
//	inv, err := r.client.Invitation(keystoreId, invitationId)
//	if err != nil {
//		return invitation.Invitation{}, errors.WithMessage(err, "failed to get invitation")
//	}
//
//	// TODO: this should be LOCAL keystoreId, while the function only receives remote one
//	keystoreKey, err := r.keystores.KeystoreKey(localKeystoreId)
//	if err != nil {
//		return invitation.Invitation{}, errors.WithMessage(err, "failed to get keystore key")
//	}
//
//	if inv.Status != "accepted" || inv.InviteeKey == nil {
//		return invitation.Invitation{}, errors.New("invitation has not been accepted")
//	}
//
//	asymKey, err := curve25519.X25519(r.privateKey, inv.InviteeKey)
//	if err != nil {
//		return invitation.Invitation{}, errors.Wrap(err, "failed to create asymmetric key")
//	}
//
//	// encrypt the keystore key with the asymmetric key
//	encryptedKeystoreKey, err := crypto.AES256CBC_Encrypt(asymKey, keystoreKey)
//	if err != nil {
//		return invitation.Invitation{}, errors.Wrap(err, "failed to encrypt keystore key")
//	}
//
//	inv, err = r.client.FinalizeInvitation(keystoreId, invitationId, encryptedKeystoreKey)
//	if err != nil {
//		return invitation.Invitation{}, errors.WithMessage(err, "failed to finalize invitation")
//	}
//
//	return inv, nil
//}
