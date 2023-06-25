package application

import (
	"github.com/pkg/errors"

	"myst/src/client/application/domain/user"
)

func (app *application) Register(sessionId []byte, username, password string) (user.User, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return user.User{}, ErrForbidden
	}

	var mustInit bool

	creds, err := app.enclave.Credentials(app.key)
	if errors.Is(err, ErrCredentialsNotFound) {
		mustInit = true
	} else if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to query credentials")
	}

	if !mustInit && creds.Address != "" && creds.Address != app.remote.Address() {
		return user.User{}, ErrRemoteAddressMismatch
	}

	u, err := app.remote.Register(username, password, creds.PublicKey)
	if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to register user")
	}

	creds.Address = app.remote.Address()
	creds.Username = username
	creds.Password = password

	_, err = app.enclave.UpdateCredentials(app.key, creds)
	if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to update credentials")
	}

	return u, nil
}

// CurrentUser returns the current user if there is one
func (app *application) CurrentUser(sessionId []byte) (*user.User, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return nil, ErrForbidden
	}

	rem, err := app.enclave.Credentials(app.key)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query credentials")
	}

	u := app.remote.CurrentUser()
	if u == nil {
		return nil, nil
	}

	u.PublicKey = rem.PublicKey

	return u, nil
}

func (app *application) SharedSecret(sessionId []byte, userId string) ([]byte, error) {
	app.mux.Lock()
	defer app.mux.Unlock()

	if !app.sessionActive(sessionId) {
		return nil, ErrForbidden
	}

	creds, err := app.enclave.Credentials(app.key)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query credentials")
	}

	return app.remote.SharedSecret(creds.PrivateKey, userId)
}
