package application

import (
	"fmt"

	"github.com/pkg/errors"

	"myst/internal/client/application/domain/enclave"
	"myst/internal/client/application/domain/entry"
	"myst/internal/client/application/domain/invitation"
	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/application/domain/user"
)

func (app *application) SignIn(username, password string) error {
	var mustInit bool

	rem, err := app.keystores.Remote()
	if errors.Is(err, enclave.ErrRemoteNotSet) {
		mustInit = true
	} else if err != nil {
		return err
	}

	if !mustInit {
		if rem.Address != app.remote.Address() {
			return fmt.Errorf("remote address mismatch")
		}

		return app.remote.SignIn(username, password)
	}

	addr := app.remote.Address()

	err = app.remote.SignIn(username, password)
	if err != nil {
		return err
	}

	return app.keystores.SetRemote(addr, username, password)
}

func (app *application) Register(username, password string) (user.User, error) {
	u, err := app.remote.Register(username, password)
	if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to register")
	}

	err = app.keystores.SetRemote(app.remote.Address(), username, password)
	if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to set user info")
	}

	err = app.remote.SignIn(username, password)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

func (app *application) CurrentUser() (*user.User, error) {
	_, err := app.keystores.Remote()
	if err != nil {
		return nil, err

	}

	return app.remote.CurrentUser(), nil
}

func (app *application) SignOut() error {
	return app.remote.SignOut()
}

func (app *application) CreateFirstKeystore(k keystore.Keystore, password string) (keystore.Keystore, error) {
	return app.keystores.CreateFirstKeystore(k, password)
}

func (app *application) CreateEnclave(password string) error {
	return app.keystores.CreateEnclave(password)
}

func (app *application) Enclave() error {
	return app.keystores.Enclave()
}

func (app *application) Remote() (enclave.Remote, error) {
	return app.keystores.Remote()
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

	var trySignIn bool
	rem, err := app.keystores.Remote()
	if err == nil {
		trySignIn = true
	} else if !errors.Is(err, enclave.ErrRemoteNotSet) {
		return err
	}

	if trySignIn {
		if rem.Address != app.remote.Address() {
			return fmt.Errorf("remote address mismatch")
		}

		err = app.remote.SignIn(rem.Username, rem.Password)
		if err != nil {
			return err
		}
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
