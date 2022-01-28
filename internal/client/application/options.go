package application

import "myst/internal/client/application/keystoreservice"

type Option func(app *application) error

func WithKeystoreRepository(repo keystoreservice.KeystoreRepository) Option {
	return func(app *application) error {
		app.repositories.keystoreRepo = repo
		return nil
	}
}
