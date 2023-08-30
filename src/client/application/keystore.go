package application

import (
	"strings"
	"time"

	"github.com/pkg/errors"

	"myst/src/client/application/domain/entry"
	"myst/src/client/application/domain/keystore"
)

func (app *application) Initialize(password string) ([]byte, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	key, err := app.enclave.Initialize(password)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to initialize enclave")
	}

	id, err := app.newSession()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create new session")
	}

	app.key = key

	return id, nil
}

func (app *application) IsInitialized(sessionId []byte) (bool, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	isInit, err := app.enclave.IsInitialized()
	if err != nil {
		return false, errors.WithMessage(err, "failed to query enclave initialization status")
	}

	if isInit && !app.sessionActive(sessionId) {
		return false, ErrAuthenticationRequired
	}

	return isInit, nil
}

func (app *application) Authenticate(password string) ([]byte, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	key, err := app.enclave.Authenticate(password)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to authenticate against enclave")
	}

	rem, err := app.enclave.Credentials(key)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query credentials")
	}

	if rem.Address != "" {
		if rem.Address != app.remote.Address() {
			return nil, ErrRemoteAddressMismatch
		}

		err = app.remote.Authenticate(rem.Username, rem.Password)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to authenticate against remote")
		}
	}

	sessionId, err := app.newSession()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create new session")
	}

	app.key = key

	return sessionId, nil
}

func (app *application) HealthCheck(sessionId []byte) error {
	app.mux.Lock()
	defer app.mux.Unlock()

	key := string(sessionId)
	app.sessions[key] = time.Now()
	return nil
}

func (app *application) CreateKeystore(sessionId []byte, name string) (keystore.Keystore, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return keystore.Keystore{}, ErrForbidden
	}

	name = strings.TrimSpace(name)

	if len(name) == 0 || len(name) > 24 {
		return keystore.Keystore{}, ErrInvalidKeystoreName
	}

	ks, err := app.enclave.Keystores(app.key)
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

	k, err = app.enclave.CreateKeystore(app.key, k)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to create keystore")
	}

	return k, nil
}

func (app *application) DeleteKeystore(sessionId []byte, id string) error {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return ErrForbidden
	}

	k, err := app.enclave.Keystore(app.key, id)
	if err != nil {
		return errors.WithMessage(err, "failed to get keystore")
	}

	if k.RemoteId != "" {
		err = app.remote.DeleteKeystore(k.RemoteId)
		if err != nil {
			return errors.WithMessage(err, "failed to delete keystore from remote")
		}
	}

	err = app.enclave.DeleteKeystore(app.key, id)
	if err != nil {
		return errors.WithMessage(err, "failed to delete keystore from enclave")
	}

	return nil
}

func (app *application) Keystore(sessionId []byte, id string) (keystore.Keystore, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return keystore.Keystore{}, ErrForbidden
	}

	k, err := app.enclave.Keystore(app.key, id)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystore")
	}

	return k, nil
}

func (app *application) keystoreByRemoteId(id string) (keystore.Keystore, error) {
	ks, err := app.enclave.Keystores(app.key)
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

func (app *application) Keystores(sessionId []byte) (map[string]keystore.Keystore, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	isInit, err := app.enclave.IsInitialized()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query enclave initialization status")
	}

	if isInit && !app.sessionActive(sessionId) {
		return nil, ErrAuthenticationRequired
	}

	ks, err := app.enclave.Keystores(app.key)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get keystores")
	}

	for id, k := range ks {
		ks[id] = k
	}

	if app.remote.Authenticated() {
		// add remote keystores too
		rem, err := app.enclave.Credentials(app.key)
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
	sessionId []byte,
	keystoreId string, website, username, password, notes string) (entry.Entry, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return entry.Entry{}, ErrForbidden
	}

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

	k, err := app.enclave.Keystore(app.key, keystoreId)
	if err != nil {
		return entry.Entry{}, errors.WithMessage(err, "failed to get keystore")
	}

	k.Entries[e.Id] = e

	_, err = app.enclave.UpdateKeystore(app.key, k)
	if err != nil {
		return entry.Entry{}, errors.WithMessage(err, "failed to update keystore")
	}

	return e, nil
}

func (app *application) UpdateEntry(
	sessionId []byte,
	keystoreId, entryId string, opts UpdateEntryOptions) (entry.Entry, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return entry.Entry{}, ErrForbidden
	}

	// do not allow empty password
	if opts.Password != nil && strings.TrimSpace(*opts.Password) == "" {
		return entry.Entry{}, ErrInvalidPassword
	}

	k, err := app.enclave.Keystore(app.key, keystoreId)
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

	_, err = app.enclave.UpdateKeystore(app.key, k)
	if err != nil {
		return entry.Entry{}, errors.WithMessage(err, "failed to update keystore")
	}

	return e, nil
}

func (app *application) DeleteEntry(sessionId []byte, keystoreId, entryId string) error {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return ErrForbidden
	}

	k, err := app.enclave.Keystore(app.key, keystoreId)
	if err != nil {
		return errors.WithMessage(err, "failed to get keystore")
	}

	if _, ok := k.Entries[entryId]; !ok {
		return ErrEntryNotFound
	}

	delete(k.Entries, entryId)

	_, err = app.enclave.UpdateKeystore(app.key, k)
	if err != nil {
		return errors.WithMessage(err, "failed to update keystore")
	}

	return nil
}
