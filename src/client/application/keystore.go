package application

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"myst/src/client/application/domain/credentials"
	"myst/src/client/application/domain/entry"
	"myst/src/client/application/domain/keystore"
)

var (
	ErrAuthenticationRequired = errors.New("authentication required")
	ErrAuthenticationFailed   = errors.New("authentication failed")
	ErrInitializationRequired = errors.New("initialization required")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrEntryNotFound          = errors.New("entry not found")
	ErrInvalidKeystoreName    = errors.New("invalid keystore name")
)

func (app *application) HealthCheck() {
	app.enclave.HealthCheck()
}

func (app *application) UpdateEntry(
	keystoreId, entryId string, opts UpdateEntryOptions) (entry.Entry, error) {
	// do not allow empty password

	if opts.Password != nil && strings.TrimSpace(*opts.Password) == "" {
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

	if opts.Password != nil {
		e.Password = *opts.Password
	}

	if opts.Notes != nil {
		e.Notes = *opts.Notes
	}

	k.Entries[e.Id] = e

	_, err = app.enclave.UpdateKeystore(k)
	if err != nil {
		return entry.Entry{}, err
	}

	return e, nil
}

func (app *application) DeleteEntry(keystoreId, entryId string) error {
	k, err := app.enclave.Keystore(keystoreId)
	if err != nil {
		return err
	}

	if _, ok := k.Entries[entryId]; !ok {
		return ErrEntryNotFound
	}

	delete(k.Entries, entryId)

	_, err = app.enclave.UpdateKeystore(k)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) CreateKeystore(name string) (keystore.Keystore, error) {
	name = strings.TrimSpace(name)

	if len(name) == 0 || len(name) > 24 {
		return keystore.Keystore{}, ErrInvalidKeystoreName
	}

	ks, err := app.enclave.Keystores()
	if err != nil {
		return keystore.Keystore{}, err
	}

	for _, k2 := range ks {
		if k2.Name == name {
			return keystore.Keystore{}, ErrInvalidKeystoreName
		}
	}

	k := keystore.New(keystore.WithName(name))

	return app.enclave.CreateKeystore(k)
}

func (app *application) DeleteKeystore(id string) error {
	k, err := app.enclave.Keystore(id)
	if err != nil {
		return err
	}

	if k.RemoteId != "" {
		err = app.remote.DeleteKeystore(k.RemoteId)
		if err != nil {
			return err
		}
	}

	return app.enclave.DeleteKeystore(id)
}

//k.Id,
//		entry.WithWebsite(req.Website),
//		entry.WithUsername(req.Username),
//		entry.WithPassword(req.Password),
//		entry.WithNotes(req.Notes),

func (app *application) CreateEntry(
	keystoreId string, website, username, password, notes string) (entry.Entry, error) {
	// do not allow empty password
	if strings.TrimSpace(password) == "" {
		return entry.Entry{}, ErrInvalidPassword
	}

	e := entry.New(
		entry.WithWebsite(website),
		entry.WithUsername(username),
		entry.WithPassword(password),
		entry.WithNotes(notes),
	)

	k, err := app.enclave.Keystore(keystoreId)
	if err != nil {
		return entry.Entry{}, err
	}

	k.Entries[e.Id] = e

	_, err = app.enclave.UpdateKeystore(k)
	if err != nil {
		return entry.Entry{}, err
	}

	return e, nil
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

// func (app *application) UpdateCredentials(creds credentials.Credentials) error {
// 	return app.enclave.UpdateCredentials(creds)
// }

func (app *application) Credentials() (credentials.Credentials, error) {
	return app.enclave.Credentials()
}
