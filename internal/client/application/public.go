package application

import (
	"errors"
	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/application/keystoreservice"
)

var (
	ErrInitializationRequired = errors.New("initialization required")
	ErrAuthenticationRequired = errors.New("authentication required")
)

func (app *application) Authenticate(password string) error {
	err := app.keystoreService.Authenticate(password)
	if err == keystoreservice.ErrAuthenticationFailed {
		return ErrAuthenticationFailed
	}

	return err
}

func (app *application) CreateKeystore(name string) (*keystore.Keystore, error) {
	return app.keystoreService.Create(name)
}

func (app *application) Initialize(name, password string) (*keystore.Keystore, error) {
	return app.keystoreService.Initialize(name, password)
}

func (app *application) Keystore(id string) (*keystore.Keystore, error) {
	return app.keystoreService.Keystore(id)
}

func (app *application) Keystores() (map[string]*keystore.Keystore, error) {
	ks, err := app.keystoreService.Keystores()
	if err == keystoreservice.ErrInitializationRequired {
		return nil, ErrInitializationRequired
	} else if err == keystoreservice.ErrAuthenticationRequired {
		return nil, ErrAuthenticationRequired
	}

	return ks, err
}

func (app *application) UpdateKeystore(k *keystore.Keystore) error {
	return app.keystoreService.Update(k)
}

func (app *application) HealthCheck() {
	app.keystoreService.HealthCheck()
}

func (app *application) SignIn(username, password string) error {
	return app.repositories.remote.SignIn(username, password)
}

func (app *application) SignOut() error {
	panic("implement me")
}
