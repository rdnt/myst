package application

import (
	"myst/internal/client/application/domain/keystore"
)

func (app *application) Authenticate(password string) error {
	err := app.keystoreService.Authenticate(password)
	//if err == keystorerepo.ErrAuthenticationFailed {
	//	return ErrAuthenticationFailed
	//}

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
	return app.keystoreService.Keystores()
}

func (app *application) KeystoreIds() ([]string, error) {
	return app.keystoreService.KeystoreIds()
}

func (app *application) UpdateKeystore(k *keystore.Keystore) error {
	return app.keystoreService.Update(k)
}

func (app *application) HealthCheck() {
	// TODO: add health check to enclavekeystorerepo
	//app.keystoreRepo.HealthCheck()
}

func (app *application) SignIn(username, password string) error {
	return app.repositories.remote.SignIn(username, password)
}

func (app *application) SignOut() error {
	panic("implement me")
}
