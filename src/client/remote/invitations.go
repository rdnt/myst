package remote

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"

	"myst/pkg/crypto"
	"myst/pkg/logger"
	"myst/src/client/application/domain/invitation"
	"myst/src/server/rest/generated"
)

var (
	ErrSignedOut = errors.New("signed out")
)

func (r *remote) Invitation(id string) (invitation.Invitation, error) {
	return r.getInvitation(id)
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
			InviterId: r.user.Id,
			InviteeId: inv.InviteeId,
		},
	)
	if err != nil {
		return invitation.Invitation{}, err
	}

	if res.JSON201 == nil {
		return invitation.Invitation{}, fmt.Errorf("invalid response")
	}

	inv, err = InvitationFromJSON(*res.JSON201)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *remote) Address() string {
	return r.address
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

func (r *remote) AcceptInvitation(invitationId string) (invitation.Invitation, error) {
	if !r.SignedIn() {
		return invitation.Invitation{}, ErrSignedOut
	}

	res, err := r.client.AcceptInvitationWithResponse(
		context.Background(), invitationId,
	)
	if err != nil {
		return invitation.Invitation{}, err
	}

	if res.JSON200 == nil {
		return invitation.Invitation{}, fmt.Errorf("invalid response")
	}

	inv, err := InvitationFromJSON(*res.JSON200)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *remote) DeclineOrCancelInvitation(id string) (invitation.Invitation, error) {
	if !r.SignedIn() {
		return invitation.Invitation{}, ErrSignedOut
	}

	res, err := r.client.DeclineOrCancelInvitationWithResponse(
		context.Background(), id,
	)
	if err != nil {
		return invitation.Invitation{}, err
	}

	if res.JSON200 == nil {
		return invitation.Invitation{}, fmt.Errorf("invalid response")
	}

	inv, err := InvitationFromJSON(*res.JSON200)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *remote) FinalizeInvitation(invitationId string, keystoreKey, privateKey []byte) (invitation.Invitation, error) {
	inv, err := r.getInvitation(invitationId)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get invitation")
	}

	if inv.Status != "accepted" {
		return invitation.Invitation{}, errors.New("invitation has not been accepted")
	}

	asymKey, err := curve25519.X25519(privateKey, inv.InviteePublicKey)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to create asymmetric key")
	}

	logger.Error("@@@ ###################### SYMMETRIC WHEN CREATE", string(asymKey))

	// encrypt the keystore key with the asymmetric key
	encryptedKeystoreKey, err := crypto.AES256CBC_Encrypt(asymKey, keystoreKey)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to encrypt keystore key")
	}

	res, err := r.client.FinalizeInvitationWithResponse(context.Background(), invitationId, generated.FinalizeInvitationJSONRequestBody{
		KeystoreKey: encryptedKeystoreKey,
	})
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to finalize invitation")
	}

	if res.JSON200 == nil {
		return invitation.Invitation{}, fmt.Errorf("invalid response")
	}

	return InvitationFromJSON(*res.JSON200)
}

func (r *remote) getInvitation(id string) (invitation.Invitation, error) {
	res, err := r.client.GetInvitationWithResponse(context.Background(), id)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get invitation")
	}

	if res.JSON200 == nil {
		return invitation.Invitation{}, fmt.Errorf("invalid response")
	}

	return InvitationFromJSON(*res.JSON200)
}

//
// func (r *remote) FinalizeInvitation(localKeystoreId, keystoreId, invitationId string) (invitation.Invitation, error) {
//	inv, err := r.client.Invitation(keystoreId, invitationId)
//	if err != nil {
//		return invitation.Invitation{}, errors.WithMessage(err, "failed to get invitation")
//	}
//
//	// TODO: this should be LOCAL keystoreId, while the function only receives remote one
//	keystoreKey, err := r.keystores.EncryptedKeystoreKey(localKeystoreId)
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
// }
