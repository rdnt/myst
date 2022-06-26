package application

import (
	"myst/internal/client/application/domain/enclave"
	"myst/internal/client/remote"
)

type Option func(app *application) error

func WithKeystoreService(service enclave.KeystoreService) Option {
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
