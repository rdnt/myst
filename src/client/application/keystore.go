package application

import (
	"strings"

	"github.com/pkg/errors"

	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/keystore/entry"
	"myst/src/client/application/domain/remote"
)

var (
	ErrInvalidKeystoreRepository = errors.New("invalid keystore repository")
	ErrAuthenticationRequired    = errors.New("authentication required")
	ErrAuthenticationFailed      = errors.New("authentication failed")
	ErrInitializationRequired    = errors.New("initialization required")
	ErrInvalidPassword           = errors.New("invalid password")
	ErrEntryNotFound             = errors.New("entry not found")
)

func (app *application) HealthCheck() {
	app.repo.HealthCheck()
}

func (app *application) UpdateKeystoreEntry(keystoreId, entryId string, password, notes *string) (entry.Entry, error) {
	// do not allow empty password
	if password != nil && strings.TrimSpace(*password) == "" {
		return entry.Entry{}, ErrInvalidPassword
	}

	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return entry.Entry{}, err
	}

	entries := k.Entries

	e, ok := entries[entryId]
	if !ok {
		return entry.Entry{}, ErrEntryNotFound
	}

	if password != nil {
		e.SetPassword(*password)
	}

	if notes != nil {
		e.SetNotes(*notes)
	}

	entries[e.Id] = e

	k.Entries = entries

	return e, app.UpdateKeystore(k)
}

func (app *application) DeleteKeystoreEntry(keystoreId, entryId string) error {
	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return err
	}

	entries := k.Entries

	if _, ok := entries[entryId]; !ok {
		return ErrEntryNotFound
	}

	delete(entries, entryId)
	k.Entries = entries

	return app.UpdateKeystore(k)
}

func (app *application) CreateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	return app.keystores.CreateKeystore(k)
}

func (app *application) DeleteKeystore(id string) error {
	return app.keystores.DeleteKeystore(id)
}

func (app *application) CreateKeystoreEntry(keystoreId string, opts ...entry.Option) (entry.Entry, error) {
	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return entry.Entry{}, err
	}

	e := entry.New(opts...)

	entries := k.Entries
	entries[e.Id] = e
	k.Entries = entries

	return e, app.keystores.UpdateKeystore(k)
}

func (app *application) Initialize(password string) error {
	return app.repo.Initialize(password)
}

func (app *application) IsInitialized() error {
	return app.repo.IsInitialized()
}

func (app *application) Keystore(id string) (keystore.Keystore, error) {
	k, err := app.keystores.Keystore(id)
	// if errors.Is(err, keystore.ErrAuthenticationRequired) {
	//	return nil, ErrAuthenticationRequired
	// }

	return k, err
}

func (app *application) KeystoreByRemoteId(id string) (keystore.Keystore, error) {
	ks, err := app.keystores.Keystores()
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystores")
	}

	for _, k2 := range ks {
		if k2.RemoteId == id {
			return k2, nil
		}
	}

	return keystore.Keystore{}, errors.New("keystore not found")
}

func (app *application) Keystores() (map[string]keystore.Keystore, error) {
	ks, err := app.keystores.Keystores()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get keystores")
	}

	for id, k := range ks {
		k.Access = "read/write"
		ks[id] = k
	}

	if app.remote.SignedIn() {
		rem, err := app.credentials.Remote()
		if err != nil {
			return nil, errors.WithMessage(err, "failed to get remote")
		}

		rks, err := app.remote.Keystores(rem.PrivateKey)
		if err != nil {
			return nil, err
		}

		for _, rk := range rks {
			if _, ok := ks[rk.Id]; !ok {
				rk.Access = "read"
				ks[rk.Id] = rk
			}
		}
	}

	return ks, nil
}

func (app *application) UpdateKeystore(k keystore.Keystore) error {
	err := app.keystores.UpdateKeystore(k)
	// if errors.Is(err, keystore.ErrAuthenticationRequired) {
	//	return ErrAuthenticationRequired
	// }

	return err
}

func (app *application) SetRemote(address, username, password string, publicKey, privateKey []byte) error {
	return app.credentials.SetRemote(address, username, password, publicKey, privateKey)
}

func (app *application) Remote() (remote.Remote, error) {
	return app.credentials.Remote()
}
