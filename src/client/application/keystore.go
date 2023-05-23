package application

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"myst/src/client/application/domain/credentials"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/keystore/entry"
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
	app.enclave.HealthCheck()
}

func (app *application) UpdateKeystoreEntry(
	keystoreId, entryId string, password, notes *string) (entry.Entry, error) {
	// do not allow empty password
	if password != nil && strings.TrimSpace(*password) == "" {
		return entry.Entry{}, ErrInvalidPassword
	}

	k, err := app.enclave.Keystore(keystoreId)
	if err != nil {
		return entry.Entry{}, err
	}

	e, ok := k.Entries[entryId]
	if !ok {
		return entry.Entry{}, ErrEntryNotFound
	}

	if password != nil {
		e.SetPassword(*password)
	}

	if notes != nil {
		e.SetNotes(*notes)
	}

	k.Entries[e.Id] = e

	return e, app.enclave.UpdateKeystore(k)
}

func (app *application) DeleteKeystoreEntry(keystoreId, entryId string) error {
	k, err := app.enclave.Keystore(keystoreId)
	if err != nil {
		return err
	}

	if _, ok := k.Entries[entryId]; !ok {
		return ErrEntryNotFound
	}

	delete(k.Entries, entryId)

	return app.enclave.UpdateKeystore(k)
}

func (app *application) CreateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	return app.enclave.CreateKeystore(k)
}

func (app *application) DeleteKeystore(id string) error {
	return app.enclave.DeleteKeystore(id)
}

func (app *application) CreateKeystoreEntry(
	keystoreId string, opts ...entry.Option) (entry.Entry, error) {
	e := entry.New(opts...)

	// do not allow empty password
	if strings.TrimSpace(e.Password) == "" {
		return entry.Entry{}, ErrInvalidPassword
	}

	k, err := app.enclave.Keystore(keystoreId)
	if err != nil {
		return entry.Entry{}, err
	}

	k.Entries[e.Id] = e

	return e, app.enclave.UpdateKeystore(k)
}

func (app *application) Initialize(password string) error {
	return app.enclave.Initialize(password)
}

func (app *application) IsInitialized() (bool, error) {
	return app.enclave.IsInitialized()
}

func (app *application) Keystore(id string) (keystore.Keystore, error) {
	return app.enclave.Keystore(id)
}

func (app *application) keystoreByRemoteId(id string) (keystore.Keystore, error) {
	ks, err := app.enclave.Keystores()
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get enclave")
	}

	for _, k2 := range ks {
		if k2.RemoteId == id {
			return k2, nil
		}
	}

	return keystore.Keystore{}, errors.New("keystore not found")
}

func (app *application) Keystores() (map[string]keystore.Keystore, error) {
	ks, err := app.enclave.Keystores()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get enclave")
	}

	for id, k := range ks {
		ks[id] = k
	}

	if app.remote.Authenticated() {
		rem, err := app.enclave.Credentials()
		if err != nil {
			return nil, errors.WithMessage(err, "failed to get remote")
		}

		rks, err := app.remote.Keystores(rem.PrivateKey)
		if err != nil {
			fmt.Print("FAILED TO GET REMOTE KEYSTORES")
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

func (app *application) UpdateCredentials(creds credentials.Credentials) error {
	return app.enclave.UpdateCredentials(creds)
}

func (app *application) Credentials() (credentials.Credentials, error) {
	return app.enclave.Credentials()
}
