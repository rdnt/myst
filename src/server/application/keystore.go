package application

import (
	"github.com/pkg/errors"

	"myst/pkg/optional"
	"myst/src/server/application/domain/invitation"
	"myst/src/server/application/domain/keystore"
)

// CreateKeystore creates a keystore with the given name, owner and payload.
// The keystore name is *not* unique in any way, it's only passed to present
// it to the invitee in case the owner decides to share the keystore with
// others. The payload is encrypted and stored as-is inside the keystore.
// If the keystore is not found, ErrKeystoreNotFound is returned.
func (app *application) CreateKeystore(name, ownerId string, payload []byte) (keystore.Keystore, error) {
	u, err := app.users.User(ownerId)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get user")
	}

	k := keystore.New(
		keystore.WithName(name),
		keystore.WithOwnerId(u.Id),
		keystore.WithPayload(payload),
	)

	k, err = app.keystores.CreateKeystore(k)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to create keystore")
	}

	return k, nil
}

// DeleteKeystore deletes a keystore with the given id. The userId initiating
// the deletion should be the owner of the keystore for it to be deleted,
// otherwise ErrForbidden is returned. Before deletion, all associated
// invitations are deleted, and then the keystore is also deleted.
// If the keystore is not found, ErrKeystoreNotFound is returned.
func (app *application) DeleteKeystore(userId string, keystoreId string) error {
	_, err := app.users.User(userId)
	if err != nil {
		return errors.WithMessage(err, "failed to get user")
	}

	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return errors.WithMessage(err, "failed to get keystore")
	}

	// check to see if the user should be able to delete this keystore
	if k.OwnerId != userId {
		invs, err := app.UserInvitations(
			userId,
			UserInvitationsOptions{Status: optional.Ref(invitation.Finalized)},
		)
		if err != nil {
			return errors.WithMessage(err, "failed to get user invitations")
		}

		var found bool
		for _, inv := range invs {
			if inv.KeystoreId == keystoreId {
				found = true
				break
			}
		}

		if found {
			return ErrForbidden
		} else {
			return ErrKeystoreNotFound
		}
	}

	invs, err := app.invitations.Invitations()
	if err != nil {
		return errors.WithMessage(err, "failed to get invitations")
	}

	// delete any associated invitations
	for _, inv := range invs {
		if inv.KeystoreId == keystoreId {
			err = app.invitations.DeleteInvitation(inv.Id)
			if err != nil {
				return errors.WithMessage(err, "failed to delete invitation")
			}
		}
	}

	err = app.keystores.DeleteKeystore(keystoreId)
	if err != nil {
		return errors.WithMessage(err, "failed to delete keystore")
	}

	return nil
}

// Keystore returns a keystore by id.
// If the keystore is not found, ErrKeystoreNotFound is returned.
func (app *application) Keystore(keystoreId string) (keystore.Keystore, error) {
	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystore")
	}

	return k, nil
}

// UserKeystore returns a keystore by id, performing the necessary checks to
// make sure the user is allowed to see it. If the user is not allowed to see
// the keystore, ErrKeystoreNotFound is returned.
// If the keystore is not found, ErrKeystoreNotFound is also returned.
func (app *application) UserKeystore(userId, keystoreId string) (keystore.Keystore, error) {
	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystore")
	}

	// check to see if the user should be able to see this keystore
	if k.OwnerId != userId {
		invs, err := app.UserInvitations(
			userId,
			UserInvitationsOptions{Status: optional.Ref(invitation.Finalized)},
		)
		if err != nil {
			return keystore.Keystore{}, errors.WithMessage(err, "failed to get user invitations")
		}

		var found bool
		for _, inv := range invs {
			if inv.KeystoreId == keystoreId {
				found = true
				break
			}
		}

		if found {
			// allow read access
			return k, nil
		} else {
			return keystore.Keystore{}, ErrKeystoreNotFound
		}
	}

	return k, nil
}

// Keystores returns all stored keystores.
func (app *application) Keystores() ([]keystore.Keystore, error) {
	ks, err := app.keystores.Keystores()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get keystores")
	}

	return ks, nil
}

// UpdateKeystore updates a keystore with the provided options. Only non-nil
// opts are processed. If the keystore is not found, ErrKeystoreNotFound is
// returned. The userId passed is the initiator of the update, and should be
// the keystore's owner, otherwise ErrForbidden is returned.
func (app *application) UpdateKeystore(userId, keystoreId string, opts UpdateKeystoreOptions) (keystore.
	Keystore,
	error) {
	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystore")
	}

	// check to see if the user should be able to see this keystore
	if k.OwnerId != userId {
		invs, err := app.UserInvitations(
			userId,
			UserInvitationsOptions{Status: optional.Ref(invitation.Finalized)},
		)
		if err != nil {
			return keystore.Keystore{}, errors.WithMessage(err, "failed to get user invitations")
		}

		var found bool
		for _, inv := range invs {
			if inv.KeystoreId == keystoreId {
				found = true
				break
			}
		}

		if found {
			return keystore.Keystore{}, ErrForbidden
		} else {
			return keystore.Keystore{}, ErrKeystoreNotFound
		}
	}

	// update the keystore's properties
	if opts.Name != nil {
		k.Name = *opts.Name
	}

	if opts.Payload != nil {
		k.Payload = *opts.Payload
	}

	k, err = app.keystores.UpdateKeystore(k)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to update keystore")
	}

	return k, nil
}

// UserKeystores returns all keystores a user is allowed to have access to.
// These include keystores where they are the owner, and also keystores they
// have been successfully invited to (any keystores granted by finalized
// invitations).
func (app *application) UserKeystores(userId string) ([]keystore.Keystore, error) {
	u, err := app.users.User(userId)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get user")
	}

	invs, err := app.UserInvitations(
		u.Id,
		UserInvitationsOptions{Status: optional.Ref(invitation.Finalized)},
	)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get user invitations")
	}

	keystores := []keystore.Keystore{}
	for _, inv := range invs {
		k, err := app.keystores.Keystore(inv.KeystoreId)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to get keystore")
		}

		keystores = append(keystores, k)
	}

	return keystores, nil
}
