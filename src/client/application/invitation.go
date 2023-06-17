package application

import (
	"github.com/pkg/errors"

	"myst/src/client/application/domain/invitation"
)

func (app *application) CreateInvitation(
	keystoreId string, inviteeUsername string) (invitation.Invitation, error) {
	k, err := app.enclave.Keystore(keystoreId)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get keystore")
	}

	// if keystore's remoteId is empty, then upload the keystore and save remoteId
	if k.RemoteId == "" {
		k, err = app.remote.CreateKeystore(k)
		if err != nil {
			return invitation.Invitation{}, errors.WithMessage(err, "failed to create keystore")
		}

		k, err = app.enclave.UpdateKeystore(k)
		if err != nil {
			return invitation.Invitation{}, errors.WithMessage(err, "failed to update keystore")
		}
	}

	inv, err := app.remote.CreateInvitation(k.RemoteId, inviteeUsername)
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

	return inv, nil
}

func (app *application) DeleteInvitation(id string) (invitation.Invitation, error) {
	inv, err := app.remote.DeleteInvitation(id)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to delete invitation")
	}

	return inv, nil
}

// FinalizeInvitation is the explicit call that the inviter makes to
// mark an invitation as finalized, calculate the shared encryption key
// using their private key and the invitee's public key, store the encryption
// key for future use (silent finalization) and finally encrypt the keystore
// key and attach it to the invitation.
func (app *application) FinalizeInvitation(invitationId, remoteKeystoreId string,
	inviteePublicKey []byte) (invitation.Invitation, error) {
	k, err := app.keystoreByRemoteId(remoteKeystoreId)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get keystore by remoteId")
	}

	creds, err := app.enclave.Credentials()
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get credentials")
	}

	inv, err := app.remote.FinalizeInvitation(invitationId, creds.PrivateKey, inviteePublicKey, k)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to finalize invitation")
	}

	return inv, nil
}

func (app *application) Invitation(id string) (invitation.Invitation, error) {
	inv, err := app.remote.Invitation(id)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get invitation")
	}

	return inv, nil
}

func (app *application) Invitations() (map[string]invitation.Invitation, error) {
	invs, err := app.remote.Invitations()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get invitations")
	}

	return invs, nil
}
