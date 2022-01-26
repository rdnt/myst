package application

import (
	"myst/internal/client/core/domain/keystore"
)

func (app *application) CreateKeystore(name string, password string) (*keystore.Keystore, error) {
	return app.keystoreService.Create(
		keystore.WithName(name),
		keystore.WithPassword(password),
	)
}

func (app *application) UnlockKeystore(keystoreId string, password string) (*keystore.Keystore, error) {
	return app.keystoreService.Unlock(keystoreId, password)
}

func (app *application) Keystore(id string) (*keystore.Keystore, error) {
	return app.keystoreService.Keystore(id)
}

func (app *application) Keystores() ([]*keystore.Keystore, error) {
	return app.keystoreService.Keystores()
}

func (app *application) KeystoreIds() ([]string, error) {
	return app.keystoreService.KeystoreIds()
}

func (app *application) UpdateKeystore(k *keystore.Keystore) error {
	return app.keystoreService.Update(k)
}

func (app *application) HealthCheck() {
	app.keystoreRepo.HealthCheck()
}

func (app *application) SignIn(username, password string) error {
	return app.remote.SignIn(username, password)
}

func (app *application) SignOut() error {
	panic("implement me")
}
