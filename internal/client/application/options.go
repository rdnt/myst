package application

import (
	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/remote"
)

type Option func(app *application) error

func WithKeystoreService(service keystore.Service) Option {
	return func(app *application) error {
		app.keystores = service
		return nil
	}
}

func WithRemote(remote remote.Remote) Option {
	return func(app *application) error {
		app.remote = remote
		return nil
	}
}
