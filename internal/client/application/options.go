package application

import (
	"myst/internal/client/application/domain/keystore"
)

type Option func(app *application) error

func WithKeystoreService(service keystore.Service) Option {
	return func(app *application) error {
		app.keystores = service
		return nil
	}
}

func WithRemoteAddress(address string) Option {
	return func(app *application) error {
		app.remoteAddress = address
		return nil
	}
}
