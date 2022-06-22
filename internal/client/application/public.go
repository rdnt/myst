package application

import (
	"github.com/pkg/errors"

	"myst/internal/client/application/domain/entry"
	"myst/internal/client/application/domain/invitation"
	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/application/domain/user"
)

func (app *application) SignIn() error {
	return app.remote.SignIn()
}

func (app *application) CurrentUser() *user.User {
	return app.remote.CurrentUser()
}

func (app *application) SignOut() error {
	return app.remote.SignOut()
}

func (app *application) CreateFirstKeystore(k keystore.Keystore, password string) (keystore.Keystore, error) {
	return app.keystores.CreateFirstKeystore(k, password)
}

func (app *application) CreateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	return app.keystores.CreateKeystore(k)
}

func (app *application) DeleteKeystore(id string) error {
	return app.keystores.DeleteKeystore(id)
}

func (app *application) Keystore(id string) (keystore.Keystore, error) {
	return app.keystores.Keystore(id)
}

func (app *application) CreateKeystoreEntry(keystoreId string, opts ...entry.Option) (entry.Entry, error) {
	return app.keystores.CreateKeystoreEntry(keystoreId, opts...)
}

func (app *application) KeystoreEntries(id string) (map[string]entry.Entry, error) {
	return app.keystores.KeystoreEntries(id)
}

func (app *application) UpdateKeystoreEntry(keystoreId string, entryId string, password, notes *string) (entry.Entry, error) {
	return app.UpdateKeystoreEntry(keystoreId, entryId, password, notes)
}

func (app *application) DeleteKeystoreEntry(keystoreId, entryId string) error {
	return app.keystores.DeleteKeystoreEntry(keystoreId, entryId)
}

func (app *application) Keystores() (map[string]keystore.Keystore, error) {
	ks, err := app.keystores.Keystores()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get keystores")
	}

	if app.remote.SignedIn() {
		rks, err := app.remote.Keystores()
		if err != nil {
			return nil, err
		}

		for _, rk := range rks {
			if _, ok := ks[rk.Id]; !ok {
				ks[rk.Id] = rk
			}
		}
	}

	return ks, nil
}

func (app *application) UpdateKeystore(k keystore.Keystore) error {
	return app.keystores.UpdateKeystore(k)
}

func (app *application) Authenticate(password string) error {
	err := app.keystores.Authenticate(password)
	if err != nil {
		return err
	}

	err = app.remote.SignIn()
	if err != nil {
		return err
	}

	return nil
}

func (app *application) CreateInvitation(keystoreId string, inviteeId string) (invitation.Invitation, error) {
	// TODO: needs refinement. app services should have access to the remote, not the other way around.
	//   for consideration: move keystoreKey to keystore.Repository (maybe the extended one)
	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get keystore")
	}

	k, err = app.remote.CreateKeystore(k)
	if err != nil {
		return invitation.Invitation{}, err
	}

	err = app.keystores.UpdateKeystore(k)
	if err != nil {
		return invitation.Invitation{}, err
	}

	inv := invitation.New(
		invitation.WithKeystoreId(k.RemoteId),
		invitation.WithInviteeId(inviteeId),
	)

	inv, err = app.remote.CreateInvitation(inv)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to create invitation")
	}

	log.Debug("invitation created", "invitation", inv)
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

	k, err := app.keystores.KeystoreByRemoteId(inv.KeystoreId)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get keystore")
	}

	inv, err = app.remote.FinalizeInvitation(id, k.Key)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to finalize invitation")
	}

	return inv, nil
}

func (app *application) Invitations() (map[string]invitation.Invitation, error) {
	return app.remote.Invitations()
}

//func (app *application) AcceptInvitation(keystoreId, invitationId string) (*invitation.Invitation, error) {
//	panic("implement me")
//	//rinv, err := app.remote.CreateInvitation(k.RemoteId(), inviteeId)
//	//if err != nil {
//	//	return nil, errors.WithMessage(err, "failed to create invitation")
//	//}
//}
