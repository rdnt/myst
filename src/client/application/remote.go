package application

import (
	"errors"
	"fmt"

	"myst/pkg/crypto"
	"myst/src/client/application/domain/enclave"
	"myst/src/client/application/domain/user"
)

func (app *application) SignIn(username, password string) error {
	var mustInit bool

	rem, err := app.keystores.Remote()
	if errors.Is(err, enclave.ErrRemoteNotSet) {
		mustInit = true
	} else if err != nil {
		return err
	}

	if !mustInit {
		if rem.Address != app.remote.Address() {
			return fmt.Errorf("remote address mismatch")
		}

		err = app.remote.SignIn(username, password)
		if err != nil {
			return err
		}
	}

	publicKey, privateKey, err := crypto.NewCurve25519Keypair()
	if err != nil {
		return err
	}

	_, err = app.remote.Register(username, password, rem.PublicKey)
	if err != nil {
		return err
	}

	err = app.keystores.SetRemote(app.remote.Address(), username, password, publicKey, privateKey)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) Register(username, password string) (user.User, error) {
	// u, err := app.remote.Register(username, password)
	// if err != nil {
	//	return user.User{}, errors.WithMessage(err, "failed to register")
	// }
	//
	// err = app.keystores.SetRemote(app.remote.Address(), username, password)
	// if err != nil {
	//	return user.User{}, errors.WithMessage(err, "failed to set user info")
	// }
	//
	// err = app.remote.SignIn(username, password)
	// if err != nil {
	//	return user.User{}, err
	// }
	//
	// return u, nil
	panic("implement me")
}

func (app *application) CurrentUser() (*user.User, error) {
	_, err := app.keystores.Remote()
	if err != nil {
		return nil, err

	}

	return app.remote.CurrentUser(), nil
}

func (app *application) SignOut() error {
	return app.remote.SignOut()
}

func (app *application) Authenticate(password string) error {
	err := app.keystores.Authenticate(password)
	if err != nil {
		return err
	}

	var trySignIn bool
	rem, err := app.keystores.Remote()
	if err == nil {
		trySignIn = true
	} else if !errors.Is(err, enclave.ErrRemoteNotSet) {
		return err
	}

	if trySignIn {
		if rem.Address != app.remote.Address() {
			return fmt.Errorf("remote address mismatch")
		}

		err = app.remote.SignIn(rem.Username, rem.Password)
		if err != nil {
			return err
		}
	}

	return nil
}
