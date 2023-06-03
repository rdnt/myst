package application

import (
	"github.com/pkg/errors"

	"myst/src/server/application/domain/invitation"
	"myst/src/server/application/domain/keystore"
)

func (app *application) CreateKeystore(name, ownerId string, payload []byte) (keystore.Keystore, error) {
	u, err := app.users.User(ownerId)
	if err != nil {
		return keystore.Keystore{}, err
	}

	k := keystore.New(
		keystore.WithName(name),
		keystore.WithOwnerId(u.Id),
		keystore.WithPayload(payload),
	)

	return app.keystores.CreateKeystore(k)
}

func (app *application) DeleteKeystore(userId string, keystoreId string) error {
	u, err := app.users.User(userId)
	if err != nil {
		return err
	}

	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return err
	}

	if k.OwnerId != u.Id {
		return errors.New("not allowed")
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

	return app.keystores.DeleteKeystore(keystoreId)
}

func (app *application) Keystore(keystoreId string) (keystore.Keystore, error) {
	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystore")
	}

	return k, nil
}

func (app *application) Keystores() ([]keystore.Keystore, error) {
	return app.keystores.Keystores()
}

func (app *application) UpdateKeystore(userId, keystoreId string, params KeystoreUpdateParams) (keystore.
	Keystore,
	error) {
	_, err := app.users.User(userId)
	if err != nil {
		return keystore.Keystore{}, err
	}

	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return keystore.Keystore{}, err
	}

	if k.OwnerId != userId {
		return keystore.Keystore{}, errors.New("not allowed")
	}

	if params.Name != nil {
		k.Name = *params.Name
	}

	if params.Payload != nil {
		k.Payload = *params.Payload
	}

	k, err = app.keystores.UpdateKeystore(k)
	if err != nil {
		return keystore.Keystore{}, err
	}

	return k, nil
}

func (app *application) UserKeystores(userId string) ([]keystore.Keystore, error) {
	u, err := app.users.User(userId)
	if err != nil {
		return nil, err
	}

	status := invitation.Finalized
	invs, err := app.UserInvitations(u.Id, UserInvitationsOptions{Status: &status})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get user invitations")
	}

	keystores := []keystore.Keystore{}
	for _, inv := range invs {
		k, err := app.keystores.Keystore(inv.KeystoreId)
		if errors.Is(err, keystore.ErrNotFound) {
			continue
		} else if err != nil {
			return nil, errors.WithMessage(err, "failed to get keystore")
		}
		keystores = append(keystores, k)
	}

	return keystores, nil
}

func (app *application) UserKeystore(userId, keystoreId string) (keystore.Keystore, error) {
	// TODO: implement
	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get user keystore")
	}

	return k, nil
}
