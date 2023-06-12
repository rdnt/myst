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
			return invitation.Invitation{}, errors.WithMessage(err, "failed to create keystore")
		}

		_, err = app.enclave.UpdateKeystore(k)
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

	// derive shared secret using the the user's private key and the invitee's public key
	sharedSecret, err := curve25519.X25519(creds.PrivateKey, inviteePublicKey)
	if err != nil {
		return invitation.Invitation{}, errors.Wrap(err, "failed to create asymmetric key")
	}

	// encrypt the keystore key with the asymmetric key
	encryptedKeystoreKey, err := crypto.AES256CBC_Encrypt(sharedSecret, k.Key)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to encrypt keystore key")
	}

	inv, err := app.remote.FinalizeInvitation(invitationId, encryptedKeystoreKey)
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

// TODO @rdnt @@@: cleanup these if not used after all
//func (app *application) sharedSecret(publicKey []byte) ([]byte, error) {
//	creds, err := app.enclave.Credentials()
//	if err != nil {
//		return nil, errors.WithMessage(err, "failed to get credentials")
//	}
//
//	sharedSecret, err := curve25519.X25519(creds.PrivateKey, publicKey)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to create shared secret")
//	}
//
//	// TODO: how to make it harder to 'decode' the icon? Generate an argon2id
//	// hash and store argon2id hash and salt to enclave for each user? And
//	// then return that instead and run it through hashicon externally?
//
//	return sharedSecret, nil
//}

//func (app *application) invitationWithIcon(inv invitation.Invitation) (invitation.Invitation, error) {
//	usr := app.remote.CurrentUser()
//	if usr == nil {
//		return invitation.Invitation{}, errors.New("no current user")
//	}
//
//	var err error
//	if inv.InviteeId != usr.Id {
//		inv.Invitee.SharedSecret, err = app.sharedSecret(inv.Invitee.PublicKey)
//		if err != nil {
//			return invitation.Invitation{}, errors.WithMessage(err, "failed to get shared secret")
//		}
//	}
//
//	if inv.Inviter.Id != usr.Id {
//		inv.Inviter.SharedSecret, err = app.sharedSecret(inv.Inviter.PublicKey)
//		if err != nil {
//			return invitation.Invitation{}, errors.WithMessage(err, "failed to get shared secret")
//		}
//	}
//
//	return inv, nil
//}
