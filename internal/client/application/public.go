package application

import (
	"github.com/pkg/errors"

	"myst/internal/client/application/domain/entry"
	"myst/internal/client/application/domain/invitation"
	"myst/internal/client/application/domain/keystore"
)

func (app *application) SignIn(username, password string) error {
	return app.remote.SignIn(username, password)
}

func (app *application) SignOut() error {
	return app.remote.SignOut()
}

func (app *application) CreateFirstKeystore(name, password string) (keystore.Keystore, error) {
	return app.keystores.CreateFirstKeystore(name, password)
}

func (app *application) CreateKeystore(name string) (keystore.Keystore, error) {
	return app.keystores.CreateKeystore(name)
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

func (app *application) KeystoreKey(keystoreId string) ([]byte, error) {
	return app.keystores.KeystoreKey(keystoreId)
}

func (app *application) Keystores() (map[string]keystore.Keystore, error) {
	return app.keystores.Keystores()
}

func (app *application) UpdateKeystore(k keystore.Keystore) error {
	return app.keystores.UpdateKeystore(k)
}

func (app *application) Authenticate(password string) error {
	return app.keystores.Authenticate(password)
}

func (app *application) CreateKeystoreInvitation(keystoreId string, inviteeId string) (*invitation.Invitation, error) {
	// TODO: needs refinement. app services should have access to the remote, not the other way around.
	//   for consideration: move keystoreKey to keystore.Repository (maybe the extended one)
	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return nil, err
	}

	// TODO: do this separately
	k, err = app.remote.CreateKeystore(k)
	if err != nil {
		return nil, err
	}

	err = app.keystores.UpdateKeystore(k)
	if err != nil {
		return nil, err
	}

	inv := invitation.New(
		invitation.WithKeystoreId(k.Id),
		invitation.WithInviteeId(inviteeId),
	)

	rinv, err := app.remote.CreateInvitation(inv)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create invitation")
	}

	log.Debug("invitation created", "invitation", rinv)
	return nil, nil
}

func (app *application) Invitations() (map[string]invitation.Invitation, error) {
	return app.remote.Invitations()
}

func (app *application) AcceptInvitation(keystoreId, invitationId string) (*invitation.Invitation, error) {
	panic("implement me")
	//rinv, err := app.remote.CreateInvitation(k.RemoteId(), inviteeId)
	//if err != nil {
	//	return nil, errors.WithMessage(err, "failed to create invitation")
	//}
}
