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

	return app.keystores.CreateKeystore(
		keystore.WithName(name),
		keystore.WithOwnerId(u.Id),
		keystore.WithPayload(payload),
	)
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

func (app *application) UpdateKeystore(userId, keystoreId string, params keystore.UpdateParams) (keystore.
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

	// if k.OwnerId != u.Id {
	// 	return keystore.Keystore{}, errors.New("not allowed")
	// }

	if params.Name != nil {
		k.Name = *params.Name
	}

	if params.Payload != nil {
		k.Payload = *params.Payload
	}

	err = app.keystores.UpdateKeystore(&k)
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
	invs, err := app.UserInvitations(u.Id, &invitation.UserInvitationsOptions{Status: &status})
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

func (app *application) UserKeystore(userId, keystoreId string) (keystore.Keystore, error) {
	// TODO: implement
	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get user keystore")
	}

	return k, nil
}
