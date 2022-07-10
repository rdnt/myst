package application

import (
	"strings"

	"github.com/pkg/errors"

	"myst/src/client/application/domain/enclave"
	"myst/src/client/application/domain/entry"
	"myst/src/client/application/domain/keystore"
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
	app.keystores.HealthCheck()
}

func (app *application) KeystoreEntries(id string) (map[string]entry.Entry, error) {
	k, err := app.keystores.Keystore(id)
	if err != nil {
		return nil, err
	}

	return k.Entries, nil
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

func (app *application) CreateFirstKeystore(k keystore.Keystore, password string) (keystore.Keystore, error) {
	err := app.keystores.CreateEnclave(password)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to initialize enclave")
	}

	k, err = app.keystores.CreateKeystore(k)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to create keystore")
	}

	return k, nil
}

func (app *application) CreateEnclave(password string) error {
	return app.keystores.CreateEnclave(password)
}

func (app *application) Enclave() error {
	return app.keystores.Enclave()
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

	if app.remote.SignedIn() {
		rem, err := app.keystores.Remote()
		if err != nil {
			return nil, errors.WithMessage(err, "failed to get remote")
		}

		rks, err := app.remote.Keystores(rem.PrivateKey)
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
	err := app.keystores.UpdateKeystore(k)
	// if errors.Is(err, keystore.ErrAuthenticationRequired) {
	//	return ErrAuthenticationRequired
	// }

	return err
}

func (app *application) SetRemote(address, username, password string, publicKey, privateKey []byte) error {
	return app.keystores.SetRemote(address, username, password, publicKey, privateKey)
}

func (app *application) Remote() (enclave.Remote, error) {
	return app.keystores.Remote()
}
