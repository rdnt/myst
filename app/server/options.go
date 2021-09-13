package server

import (
	"myst/app/server/domain/keystore"
	"myst/app/server/domain/user"
)

type Option func(app *Application) error

func WithKeystoreRepository(repo keystore.Repository) Option {
	return func(app *Application) error {
		app.keystorerepo = repo
		return nil
	}
}

func WithUserRepository(repo user.Repository) Option {
	return func(app *Application) error {
		app.userrepo = repo
		return nil
	}
}

func WithUserService(service user.Service) Option {
	return func(app *Application) error {
		app.userService = service
		return nil
	}
}
