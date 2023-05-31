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

	return app.invitationWithIcon(inv)
}

func (app *application) AcceptInvitation(id string) (invitation.Invitation, error) {
	inv, err := app.remote.AcceptInvitation(id)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to accept invitation")
	}

	return app.invitationWithIcon(inv)
}

func (app *application) DeclineOrCancelInvitation(id string) (invitation.Invitation, error) {
	inv, err := app.remote.DeclineOrCancelInvitation(id)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to decline or cancel invitation")
	}

	return app.invitationWithIcon(inv)
}

// VerifyInvitation is called by the invitee, so that they can validate the
// inviter's keys are trustworthy and store them locally.
// If the user chooses to create an invitation towards the former inviter, now
// invitee, the invitation will be silently finalized by the client and
// the encryption key will be reused.
// func (app *application) VerifyInvitation(id string) error {
// 	inv, err := app.remote.Invitation(id)
// 	if err != nil {
// 		return errors.WithMessage(err, "failed to get invitation")
// 	}
//
// 	creds, err := app.enclave.Credentials()
// 	if err != nil {
// 		return errors.WithMessage(err, "failed to get credentials")
// 	}
//
// 	if inv.Status != invitation.Accepted && inv.Status != invitation.Finalized {
// 		return errors.New("invitation has not been accepted or finalized")
// 	}
//
// 	// derive asymmetric key using the invitee's private key and the inviter's public key
// 	asymKey, err := curve25519.X25519(creds.PrivateKey, inv.Inviter.PublicKey)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to create asymmetric key")
// 	}
//
// 	if existing, ok := creds.UserKeys[inv.Inviter.Id]; ok {
// 		if subtle.ConstantTimeCompare(existing, asymKey) != 0 {
// 			return errors.New("inviter's public key has been altered")
// 		}
// 	}
//
// 	creds.UserKeys[inv.Inviter.Id] = asymKey
//
// 	err = app.enclave.UpdateCredentials(creds)
// 	if err != nil {
// 		return errors.WithMessage(err, "failed to update credentials")
// 	}
//
// 	return nil
// }

func (app *application) sharedSecret(publicKey []byte) ([]byte, error) {
	creds, err := app.enclave.Credentials()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get credentials")
	}

	asymKey, err := curve25519.X25519(creds.PrivateKey, publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create asymmetric key")
	}

	// TODO: how to make it harder to 'decode' the icon? Generate an argon2id
	// hash and store argon2id hash and salt to enclave for each user? And
	// then return that instead and run it through hashicon externally?

	return asymKey, nil
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
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get keystore")
	}

	creds, err := app.enclave.Credentials()
	if err != nil {
		return invitation.Invitation{}, err
	}

	// derive asymmetric key using the inviter's private key and the invitee's public key
	asymKey, err := curve25519.X25519(creds.PrivateKey, inviteePublicKey)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to create asymmetric key")
	}

	// if existing, ok := creds.UserKeys[inv.Invitee.Id]; ok {
	// 	if subtle.ConstantTimeCompare(existing, asymKey) != 0 {
	// 		return invitation.Invitation{}, errors.New("invitee's public key has been altered")
	// 	}
	// }

	// creds.UserKeys[inv.Invitee.Id] = asymKey

	// err = app.enclave.UpdateCredentials(creds)
	// if err != nil {
	// 	return invitation.Invitation{}, errors.WithMessage(err, "failed to update credentials")
	// }

	// encrypt the keystore key with the asymmetric key
	encryptedKeystoreKey, err := crypto.AES256CBC_Encrypt(asymKey, k.Key)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to encrypt keystore key")
	}

	inv, err := app.remote.FinalizeInvitation(invitationId, encryptedKeystoreKey)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to finalize invitation")
	}

	return app.invitationWithIcon(inv)
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

	// creds, err := app.enclave.Credentials()
	// if err != nil {
	// 	return invitation.Invitation{}, errors.WithMessage(err, "failed to get credentials")
	// }
	//
	// if _, ok := creds.UserKeys[inv.Inviter.Id]; ok && usr.Id != inv.Inviter.Id {
	// 	inv.Inviter.Verified = true
	// }
	//
	// if _, ok := creds.UserKeys[inv.Invitee.Id]; ok && usr.Id != inv.Invitee.Id {
	// 	inv.Invitee.Verified = true
	// }

	return app.invitationWithIcon(inv)
}

func (app *application) Invitations() (map[string]invitation.Invitation, error) {
	usr := app.remote.CurrentUser()
	if usr == nil {
		return nil, errors.New("no current user")
	}

	invs, err := app.remote.Invitations()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get invitations")
	}

	for id, inv := range invs {
		inv, err := app.invitationWithIcon(inv)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to get invitation with icon")
		}

		invs[id] = inv
	}

	return invs, nil
}

func (app *application) invitationWithIcon(inv invitation.Invitation) (invitation.Invitation, error) {
	usr := app.remote.CurrentUser()
	if usr == nil {
		return invitation.Invitation{}, errors.New("no current user")
	}

	var err error
	if inv.Invitee.Id != usr.Id {
		inv.Invitee.SharedSecret, err = app.sharedSecret(inv.Invitee.PublicKey)
		if err != nil {
			return invitation.Invitation{}, errors.WithMessage(err, "failed to get shared secret")
		}
	}

	if inv.Inviter.Id != usr.Id {
		inv.Inviter.SharedSecret, err = app.sharedSecret(inv.Inviter.PublicKey)
		if err != nil {
			return invitation.Invitation{}, errors.WithMessage(err, "failed to get shared secret")
		}
	}

	return inv, nil
}

// func (app *application) AcceptInvitation(keystoreId, invitationId string) (*invitation.Invitation, error) {
//	panic("implement me")
//	//rinv, err := app.remote.CreateInvitation(k.RemoteId(), inviteeId)
//	//if err != nil {
//	//	return nil, errors.WithMessage(err, "failed to create invitation")
//	//}
// }
