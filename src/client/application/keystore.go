package application

import (
	"strings"

	"github.com/pkg/errors"

	"myst/src/client/application/domain/entry"
	"myst/src/client/application/domain/keystore"
)

var (
	ErrKeystoreNotFound       = errors.New("keystore not found")
	ErrAuthenticationRequired = errors.New("authentication required")
	ErrAuthenticationFailed   = errors.New("authentication failed")
	ErrInitializationRequired = errors.New("initialization required")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrInvalidWebsite         = errors.New("invalid website")
	ErrorInvalidUsername      = errors.New("invalid username")
	ErrEntryNotFound          = errors.New("entry not found")
	ErrInvalidKeystoreName    = errors.New("invalid keystore name")
	ErrCredentialsNotFound    = errors.New("credentials not found")
)

func (app *application) Initialize(password string) error {
	err := app.enclave.Initialize(password)
	if err != nil {
		return errors.WithMessage(err, "failed to initialize enclave")
	}

	return nil
}

func (app *application) IsInitialized() (bool, error) {
	isInit, err := app.enclave.IsInitialized()
	if err != nil {
		return false, errors.WithMessage(err, "failed to query enclave initialization status")
	}

	return isInit, nil
}

func (app *application) Authenticate(password string) error {
	err := app.enclave.Authenticate(password)
	if err != nil {
		return errors.WithMessage(err, "failed to authenticate against enclave")
	}

	var trySignIn bool
	rem, err := app.enclave.Credentials()
	if err == nil {
		trySignIn = true
	} else if !errors.Is(err, ErrCredentialsNotFound) {
		return errors.WithMessage(err, "failed to query credentials")
	}

	if trySignIn {
		if rem.Address != app.remote.Address() {
			return ErrRemoteAddressMismatch
		}

		// TODO: do this in a goroutine on interval to keep JWT fresh
		err = app.remote.Authenticate(rem.Username, rem.Password)
		if err != nil {
			return errors.WithMessage(err, "failed to authenticate against remote")
		}
	}

	return nil
}

func (app *application) HealthCheck() {
	app.enclave.HealthCheck()
}

func (app *application) CreateKeystore(name string) (keystore.Keystore, error) {
	name = strings.TrimSpace(name)

	if len(name) == 0 || len(name) > 24 {
		return keystore.Keystore{}, ErrInvalidKeystoreName
	}

	ks, err := app.enclave.Keystores()
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to query keystores")
	}

	// do not allow keystores with the same name
	for _, k2 := range ks {
		if k2.Name == name {
			return keystore.Keystore{}, ErrInvalidKeystoreName
		}
	}

	k := keystore.New(keystore.WithName(name))

	k, err = app.enclave.CreateKeystore(k)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to create keystore")
	}

	return k, nil
}

func (app *application) DeleteKeystore(id string) error {
	k, err := app.enclave.Keystore(id)
	if err != nil {
		return errors.WithMessage(err, "failed to get keystore")
	}

	if k.RemoteId != "" {
		err = app.remote.DeleteKeystore(k.RemoteId)
		if err != nil {
			return errors.WithMessage(err, "failed to delete keystore from remote")
		}
	}

	err = app.enclave.DeleteKeystore(id)
	if err != nil {
		return errors.WithMessage(err, "failed to delete keystore from enclave")
	}

	return nil
}

func (app *application) Keystore(id string) (keystore.Keystore, error) {
	k, err := app.enclave.Keystore(id)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystore")
	}

	return k, nil
}

func (app *application) keystoreByRemoteId(id string) (keystore.Keystore, error) {
	ks, err := app.enclave.Keystores()
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystores")
	}

	for _, k2 := range ks {
		if k2.RemoteId == id {
			return k2, nil
		}
	}

	return keystore.Keystore{}, ErrKeystoreNotFound
}

func (app *application) Keystores() (map[string]keystore.Keystore, error) {
	ks, err := app.enclave.Keystores()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get keystores")
	}

	for id, k := range ks {
		ks[id] = k
	}

	if app.remote.Authenticated() {
		// add remote keystores too

		rem, err := app.enclave.Credentials()
		if err != nil {
			return nil, errors.WithMessage(err, "failed to get credentials")
		}

		rks, err := app.remote.Keystores(rem.PrivateKey)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to get remote keystores")
		}

		for _, rk := range rks {
			if _, ok := ks[rk.Id]; !ok {
				ks[rk.Id] = rk
			}
		}
	}

	return ks, nil
}

func (app *application) CreateEntry(
	keystoreId string, website, username, password, notes string) (entry.Entry, error) {
	// do not allow empty fields for website, username, password
	if strings.TrimSpace(website) == "" {
		return entry.Entry{}, ErrInvalidWebsite
	}

	if strings.TrimSpace(username) == "" {
		return entry.Entry{}, ErrorInvalidUsername
	}

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
		return entry.Entry{}, errors.WithMessage(err, "failed to get keystore")
	}

	k.Entries[e.Id] = e

	_, err = app.enclave.UpdateKeystore(k)
	if err != nil {
		return entry.Entry{}, errors.WithMessage(err, "failed to update keystore")
	}

	return e, nil
}

func (app *application) UpdateEntry(
	keystoreId, entryId string, opts UpdateEntryOptions) (entry.Entry, error) {
	// do not allow empty password
	if opts.Password != nil && strings.TrimSpace(*opts.Password) == "" {
		return entry.Entry{}, ErrInvalidPassword
	}

	k, err := app.enclave.Keystore(keystoreId)
	if err != nil {
		return entry.Entry{}, errors.WithMessage(err, "failed to get keystore")
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
		return entry.Entry{}, errors.WithMessage(err, "failed to update keystore")
	}

	return e, nil
}

func (app *application) DeleteEntry(keystoreId, entryId string) error {
	k, err := app.enclave.Keystore(keystoreId)
	if err != nil {
		return errors.WithMessage(err, "failed to get keystore")
	}

	if _, ok := k.Entries[entryId]; !ok {
		return ErrEntryNotFound
	}

	delete(k.Entries, entryId)

	_, err = app.enclave.UpdateKeystore(k)
	if err != nil {
		return errors.WithMessage(err, "failed to update keystore")
	}

	return nil
}
