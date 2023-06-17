package application

import (
	"github.com/pkg/errors"

	"myst/src/client/application/domain/user"
)

func (app *application) Register(username, password string) (user.User, error) {
	var mustInit bool

	creds, err := app.enclave.Credentials()
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

	_, err = app.enclave.UpdateCredentials(creds)
	if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to update credentials")
	}

	return u, nil
}

// CurrentUser returns the current user if there is one
func (app *application) CurrentUser() (*user.User, error) {
	rem, err := app.enclave.Credentials()
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

func (app *application) SharedSecret(userId string) ([]byte, error) {
	creds, err := app.enclave.Credentials()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to query credentials")
	}

	return app.remote.SharedSecret(creds.PrivateKey, userId)
}
