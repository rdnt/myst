package remote

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"

	"myst/pkg/crypto"
	"myst/src/client/application"
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/server/rest/generated"
)

func (r *Remote) Invitation(id string) (invitation.Invitation, error) {
	if !r.Authenticated() {
		return invitation.Invitation{}, application.ErrAuthenticationRequired
	}

	res, err := r.client.GetInvitationWithResponse(context.Background(), id)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to get invitation")
	}

	if res.StatusCode() == http.StatusNotFound {
		return invitation.Invitation{}, application.ErrInvitationNotFound
	}

	if res.StatusCode() != http.StatusOK || res.JSON200 == nil {
		return invitation.Invitation{}, errors.New("invalid response")
	}

	inv, err := InvitationFromJSON(*res.JSON200)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *Remote) CreateInvitation(keystoreRemoteId, inviteeUsername string) (invitation.Invitation, error) {
	if !r.Authenticated() {
		return invitation.Invitation{}, application.ErrAuthenticationRequired
	}

	res, err := r.client.CreateInvitationWithResponse(
		context.Background(), keystoreRemoteId, generated.CreateInvitationJSONRequestBody{
			Invitee: inviteeUsername,
		},
	)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to create invitation")
	}

	if res.StatusCode() != http.StatusCreated || res.JSON201 == nil {
		return invitation.Invitation{}, errors.New("invalid response")
	}

	inv, err := InvitationFromJSON(*res.JSON201)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *Remote) Invitations() (map[string]invitation.Invitation, error) {
	if !r.Authenticated() {
		return nil, application.ErrAuthenticationRequired
	}

	res, err := r.client.InvitationsWithResponse(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get invitations")
	}

	if res.StatusCode() != http.StatusOK || res.JSON200 == nil {
		return nil, errors.New("invalid response")
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

func (r *Remote) AcceptInvitation(invitationId string) (invitation.Invitation, error) {
	res, err := r.client.AcceptInvitationWithResponse(
		context.Background(), invitationId,
	)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to accept invitation")
	}

	if res.StatusCode() == http.StatusNotFound {
		return invitation.Invitation{}, application.ErrInvitationNotFound
	} else if res.StatusCode() == http.StatusForbidden {
		return invitation.Invitation{}, application.ErrForbidden
	} else if res.StatusCode() != http.StatusOK || res.JSON200 == nil {
		return invitation.Invitation{}, errors.New("invalid response")
	}

	inv, err := InvitationFromJSON(*res.JSON200)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *Remote) DeleteInvitation(id string) (invitation.Invitation, error) {
	if !r.Authenticated() {
		return invitation.Invitation{}, application.ErrAuthenticationRequired
	}

	res, err := r.client.DeleteInvitationWithResponse(
		context.Background(), id,
	)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to delete invitation")
	}

	if res.StatusCode() != http.StatusOK || res.JSON200 == nil {
		return invitation.Invitation{}, errors.New("invalid response")
	}

	inv, err := InvitationFromJSON(*res.JSON200)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}

func (r *Remote) FinalizeInvitation(
	invitationId string, privateKey []byte, inviteePublicKey []byte, k keystore.Keystore,
) (invitation.Invitation, error) {
	// derive shared secret using the user's private key and the invitee's public key
	sharedSecret, err := curve25519.X25519(privateKey, inviteePublicKey)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to derive shared secret")
	}

	// encrypt the keystore key with the asymmetric key
	encryptedKeystoreKey, err := crypto.AES256CBC_Encrypt(sharedSecret, k.Key)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to encrypt keystore key")
	}

	res, err := r.client.FinalizeInvitationWithResponse(context.Background(),
		invitationId, generated.FinalizeInvitationJSONRequestBody{
			KeystoreKey: encryptedKeystoreKey,
		},
	)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to finalize invitation")
	}

	if res.StatusCode() != http.StatusOK || res.JSON200 == nil {
		return invitation.Invitation{}, errors.New("invalid response")
	}

	inv, err := InvitationFromJSON(*res.JSON200)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to parse invitation")
	}

	return inv, nil
}
