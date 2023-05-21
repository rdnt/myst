package application

import (
	"github.com/pkg/errors"

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
		return user.User{}, errors.Wrap(err, "failed to query credentials")
	}

	if !mustInit && rem.Address != app.remote.Address() {
		return user.User{}, ErrRemoteAddressMismatch
	}

	publicKey, privateKey, err := crypto.NewCurve25519Keypair()
	if err != nil {
		return user.User{}, errors.Wrap(err, "failed to generate keypair")
	}

	u, err := app.remote.Register(username, password, publicKey)
	if err != nil {
		return user.User{}, errors.Wrap(err, "failed to register user")
	}

	err = app.enclave.SetCredentials(app.remote.Address(), username, password, publicKey, privateKey)
	if err != nil {
		return user.User{}, errors.Wrap(err, "failed to update credentials")
	}

	return u, nil
}

// SignIn signs the user in against the remote. Username and password is used
// for authentication, and the publicKey is used to replace the upstream
// public key in order to be able to sync keystores.
// func (app *application) SignIn(username, password string) (user.User, error) {
// 	var mustInit bool
//
// 	rem, err := app.enclave.Credentials()
// 	if errors.Is(err, enclave.ErrRemoteNotSet) {
// 		mustInit = true
// 	} else if err != nil {
// 		return user.User{}, err
// 	}
//
// 	if !mustInit {
// 		if rem.Address != app.remote.Address() {
// 			return user.User{}, ErrRemoteAddressMismatch
// 		}
// 	}
//
// 	publicKey, privateKey, err := crypto.NewCurve25519Keypair()
// 	if err != nil {
// 		return user.User{}, err
// 	}
//
// 	u, err := app.remote.Authenticate(username, password, publicKey)
// 	if err != nil {
// 		return user.User{}, err
// 	}
//
// 	// update credentials
// 	err = app.enclave.SetCredentials(app.remote.Address(), username, password, publicKey, privateKey)
// 	if err != nil {
// 		return user.User{}, err
// 	}
//
// 	return u, nil
// }

// func (app *application) SignOut() error {
// 	rem, err := app.enclave.Credentials()
// 	if err != nil {
// 		return err
// 	}
//
// 	if rem.Address != app.remote.Address() {
// 		return ErrRemoteAddressMismatch
// 	}
//
// 	err = app.remote.SignOut()
// 	if err != nil {
// 		return err
// 	}
//
// 	// update credentials, keeping addr, pub and priv keys intact
// 	err = app.enclave.SetCredentials(app.remote.Address(), "", "", rem.PublicKey, rem.PrivateKey)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// CurrentUser returns the current user if there is one
func (app *application) CurrentUser() (*user.User, error) {
	rem, err := app.enclave.Credentials()
	if err != nil {
		return nil, err
	}

	u := app.remote.CurrentUser()
	if u == nil {
		return nil, nil
	}
	u.PublicKey = rem.PublicKey

	return u, nil
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
			return ErrRemoteAddressMismatch
		}

		err = app.remote.Authenticate(rem.Username, rem.Password)
		if err != nil {
			return err
		}
	}

	return nil
}
