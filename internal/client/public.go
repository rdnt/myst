package application

import (
	"myst/internal/client/core/domain/keystore"
)

func (app *application) CreateKeystore(name string, passphrase string) (*keystore.Keystore, error) {
	return app.keystoreService.Create(
		keystore.WithName(name),
		keystore.WithPassphrase(passphrase),
	)
}

func (app *application) UnlockKeystore(keystoreId string, passphrase string) (*keystore.Keystore, error) {
	return app.keystoreService.Unlock(keystoreId, passphrase)
}

func (app *application) Keystore(id string) (*keystore.Keystore, error) {
	return app.keystoreService.Keystore(id)
}

func (app *application) UpdateKeystore(k *keystore.Keystore) error {
	return app.keystoreService.Update(k)
}

func (app *application) HealthCheck() {
	app.keystoreRepo.HealthCheck()
}

func (app *application) SignIn(username, password string) error {
	panic("implement me")
}

func (app *application) SignOut() error {
	panic("implement me")
}
