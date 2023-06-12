package application

import (
	"github.com/pkg/errors"

	"myst/src/client/application/domain/credentials"
	"myst/src/client/enclaverepo/enclave"

	"myst/pkg/crypto"
	"myst/src/client/application/domain/user"
)

var ErrRemoteAddressMismatch = errors.New("remote address mismatch")

func (app *application) Register(username, password string) (user.User, error) {
	var mustInit bool

	rem, err := app.enclave.Credentials()
	if errors.Is(err, enclave.ErrRemoteNotSet) {
		mustInit = true
	} else if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to query credentials")
	}

	if !mustInit && rem.Address != app.remote.Address() {
		return user.User{}, ErrRemoteAddressMismatch
	}

	publicKey, privateKey, err := crypto.NewCurve25519Keypair()
	if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to generate keypair")
	}

	u, err := app.remote.Register(username, password, publicKey)
	if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to register user")
	}

	_, err = app.enclave.UpdateCredentials(credentials.Credentials{
		Address:    app.remote.Address(),
		Username:   username,
		Password:   password,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	})
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
