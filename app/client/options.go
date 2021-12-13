package client

import (
	"myst/app/client/core/domain/keystore"
)

type Option func(app *application) error

func WithKeystoreRepository(repo keystore.Repository) Option {
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
