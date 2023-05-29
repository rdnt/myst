package application

import (
	"crypto/subtle"

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

	// if keystore's remoteId is empty, then upload the keystore and save remoteId
	if k.RemoteId == "" {
		k, err = app.remote.CreateKeystore(k)
		if err != nil {
			return invitation.Invitation{}, err
		}

		err = app.enclave.UpdateKeystore(k)
		if err != nil {
			return invitation.Invitation{}, err
		}
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

// VerifyInvitation is called by the invitee, so that they can validate the
// inviter's keys are trustworthy and store them locally.
// If the user chooses to create an invitation towards the former inviter, now
// invitee, the invitation will be silently finalized by the client and
// the encryption key will be reused.
func (app *application) VerifyInvitation(id string) error {
	inv, err := app.remote.Invitation(id)
	if err != nil {
		return errors.WithMessage(err, "failed to get invitation")
	}

	creds, err := app.enclave.Credentials()
	if err != nil {
		return errors.WithMessage(err, "failed to get credentials")
	}

	if inv.Status != invitation.Accepted && inv.Status != invitation.Finalized {
		return errors.New("invitation has not been accepted or finalized")
	}

	// derive asymmetric key using the invitee's private key and the inviter's public key
	asymKey, err := curve25519.X25519(creds.PrivateKey, inv.Inviter.PublicKey)
	if err != nil {
		return errors.Wrap(err, "failed to create asymmetric key")
	}

	if existing, ok := creds.UserKeys[inv.Inviter.Id]; ok {
		if subtle.ConstantTimeCompare(existing, asymKey) != 0 {
			return errors.New("inviter's public key has been altered")
		}
	}

	creds.UserKeys[inv.Inviter.Id] = asymKey

	err = app.enclave.UpdateCredentials(creds)
	if err != nil {
		return errors.WithMessage(err, "failed to update credentials")
	}

	return nil
}

// FinalizeInvitation is the explicit call that the inviter makes to
// mark an invitation as finalized, calculate the shared encryption key
// using their private key and the invitee's public key, store the encryption
// key for future use (silent finalization) and finally encrypt the keystore
// key and attach it to the invitation.
func (app *application) FinalizeInvitation(id string) (invitation.Invitation, error) {
	inv, err := app.remote.Invitation(id)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get invitation")
	}

	creds, err := app.enclave.Credentials()
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

	// derive asymmetric key using the inviter's private key and the invitee's public key
	asymKey, err := curve25519.X25519(creds.PrivateKey, inv.Invitee.PublicKey)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to create asymmetric key")
	}

	if existing, ok := creds.UserKeys[inv.Invitee.Id]; ok {
		if subtle.ConstantTimeCompare(existing, asymKey) != 0 {
			return invitation.Invitation{}, errors.New("invitee's public key has been altered")
		}
	}

	creds.UserKeys[inv.Invitee.Id] = asymKey

	err = app.enclave.UpdateCredentials(creds)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to update credentials")
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

func (app *application) Invitation(id string) (invitation.Invitation, error) {
	usr := app.remote.CurrentUser()
	if usr == nil {
		return invitation.Invitation{}, errors.New("no current user")
	}

	inv, err := app.remote.Invitation(id)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get invitation")
	}

	creds, err := app.enclave.Credentials()
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get credentials")
	}

	if _, ok := creds.UserKeys[inv.Inviter.Id]; ok && usr.Id != inv.Inviter.Id {
		inv.Inviter.Verified = true
	}

	if _, ok := creds.UserKeys[inv.Invitee.Id]; ok && usr.Id != inv.Invitee.Id {
		inv.Invitee.Verified = true
	}

	return inv, nil
}

func (app *application) Invitations() (map[string]invitation.Invitation, error) {
	usr := app.remote.CurrentUser()
	if usr == nil {
		return nil, errors.New("no current user")
	}

	return app.remote.Invitations()
}

// func (app *applicationrefactor) AcceptInvitation(keystoreId, invitationId string) (*invitation.Invitation, error) {
//	panic("implement me")
//	//rinv, err := app.remote.CreateInvitation(k.RemoteId(), inviteeId)
//	//if err != nil {
//	//	return nil, errors.WithMessage(err, "failed to create invitation")
//	//}
// }
