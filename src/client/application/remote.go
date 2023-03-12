package application

import (
	"errors"
	"fmt"

	"myst/pkg/crypto"
	"myst/src/client/application/domain/enclave"
	"myst/src/client/application/domain/user"
)

func (app *application) Register(username, password string) (user.User, error) {
	var mustInit bool

	rem, err := app.enclave.Credentials()
	if errors.Is(err, enclave.ErrRemoteNotSet) {
		mustInit = true
	} else if err != nil {
		return user.User{}, err
	}

	if !mustInit {
		if rem.Address != app.remote.Address() {
			return user.User{}, fmt.Errorf("remote address mismatch")
		}

		_, err = app.remote.SignIn(username, password, rem.PublicKey)
		if err != nil {
			return user.User{}, err
		}
	}

	publicKey, privateKey, err := crypto.NewCurve25519Keypair()
	if err != nil {
		return user.User{}, err
	}

	u, err := app.remote.Register(username, password, publicKey)
	if err != nil {
		return user.User{}, err
	}

	err = app.enclave.SetCredentials(app.remote.Address(), username, password, publicKey, privateKey)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

func (app *application) SignIn(username, password string) (user.User, error) {
	panic("disabled")

	var mustInit bool

	rem, err := app.enclave.Credentials()
	if errors.Is(err, enclave.ErrRemoteNotSet) {
		mustInit = true
	} else if err != nil {
		return user.User{}, err
	}

	if !mustInit {
		if rem.Address != app.remote.Address() {
			return user.User{}, fmt.Errorf("remote address mismatch")
		}

		_, err = app.remote.SignIn(username, password, rem.PublicKey)
		if err != nil {
			return user.User{}, err
		}
	}

	publicKey, privateKey, err := crypto.NewCurve25519Keypair()
	if err != nil {
		return user.User{}, err
	}

	u, err := app.remote.SignIn(username, password, rem.PublicKey)
	if err != nil {
		return user.User{}, err
	}

	err = app.enclave.SetCredentials(app.remote.Address(), username, password, publicKey, privateKey)
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}

func (app *application) CurrentUser() (*user.User, error) {
	rem, err := app.enclave.Credentials()
	if err != nil {
		return nil, err

	}

	u := app.remote.CurrentUser()
	u.PublicKey = rem.PublicKey

	return u, nil
}

func (app *application) SignOut() error {
	return app.remote.SignOut()
}

func (app *application) Authenticate(password string) error {
	err := app.enclave.Authenticate(password)
	if err != nil {
		return err
	}

	var trySignIn bool
	rem, err := app.enclave.Credentials()
	if err == nil {
		trySignIn = true
	} else if !errors.Is(err, enclave.ErrRemoteNotSet) {
		return err
	}

	if trySignIn {
		if rem.Address != app.remote.Address() {
			return fmt.Errorf("remote address mismatch")
		}

		_, err = app.remote.SignIn(rem.Username, rem.Password, rem.PublicKey)
		if err != nil {
			return err
		}
	}

	return nil
}
