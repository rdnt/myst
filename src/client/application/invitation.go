package application

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"

	"myst/pkg/crypto"
	"myst/src/client/application/domain/invitation"
)

func (app *application) CreateInvitation(
	keystoreId string, inviteeUsername string) (invitation.Invitation, error) {
	k, err := app.enclave.Keystore(keystoreId)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get keystore")
	}

	k, err = app.remote.CreateKeystore(k)
	if err != nil {
		return invitation.Invitation{}, err
	}

	err = app.enclave.UpdateKeystore(k)
	if err != nil {
		return invitation.Invitation{}, err
	}

	inv := invitation.New(
		invitation.WithKeystore(k),
		invitation.WithInviteeUsername(inviteeUsername),
	)

	inv, err = app.remote.CreateInvitation(inv)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to create invitation")
	}

	return inv, nil
}

func (app *application) AcceptInvitation(id string) (invitation.Invitation, error) {
	inv, err := app.remote.AcceptInvitation(id)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to accept invitation")
	}

	return inv, err
}

func (app *application) DeclineOrCancelInvitation(id string) (invitation.Invitation, error) {
	inv, err := app.remote.DeclineOrCancelInvitation(id)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to decline or cancel invitation")
	}

	return inv, err
}

func (app *application) FinalizeInvitation(id string) (invitation.Invitation, error) {
	inv, err := app.remote.Invitation(id)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get invitation")
	}

	rem, err := app.enclave.Credentials()
	if err != nil {
		return invitation.Invitation{}, err
	}

	k, err := app.keystoreByRemoteId(inv.Keystore.RemoteId)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get keystore")
	}

	if inv.Status != "accepted" {
		return invitation.Invitation{}, errors.New("invitation has not been accepted")
	}

	asymKey, err := curve25519.X25519(rem.PrivateKey, inv.Invitee.PublicKey)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to create asymmetric key")
	}

	// encrypt the keystore key with the asymmetric key
	encryptedKeystoreKey, err := crypto.AES256CBC_Encrypt(asymKey, k.Key)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to encrypt keystore key")
	}

	inv, err = app.remote.FinalizeInvitation(id, encryptedKeystoreKey)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to finalize invitation")
	}

	return inv, nil
}

func (app *application) Invitations() (map[string]invitation.Invitation, error) {
	return app.remote.Invitations()
}

func (app *application) Invitation(id string) (invitation.Invitation, error) {
	return app.remote.Invitation(id)
}

// func (app *applicationrefactor) AcceptInvitation(keystoreId, invitationId string) (*invitation.Invitation, error) {
//	panic("implement me")
//	//rinv, err := app.remote.CreateInvitation(k.RemoteId(), inviteeId)
//	//if err != nil {
//	//	return nil, errors.WithMessage(err, "failed to create invitation")
//	//}
// }
