package client

import (
	"myst/internal/client/core/domain/keystore"
)

type Option func(app *application) error

func WithKeystoreRepository(repo KeystoreRepository) Option {
	return func(app *application) error {
		app.keystoreRepo = repo
		return nil
	}
}

func WithKeystoreService(service keystore.Service) Option {
	return func(app *application) error {
		app.keystoreService = service
		return nil
	}
}
