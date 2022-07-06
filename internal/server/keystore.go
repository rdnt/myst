package application

import (
	"github.com/pkg/errors"

	"myst/internal/server/core/domain/invitation"
	"myst/internal/server/core/domain/keystore"
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
	panic("implement me")
	// k, err := app.keystores.UserKeystore(userId, keystoreId)
	// if err != nil {
	// 	return keystore.Keystore{}, errors.WithMessage(err, "failed to get user keystore")
	// }
	//
	// return k, nil
}
