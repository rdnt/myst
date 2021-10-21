package core

import (
	"myst/app/server/core/domain/invitation"
	"myst/app/server/core/domain/keystore"
	"myst/app/server/core/domain/user"
)

type Option func(app *Application) error

func WithKeystoreRepository(repo keystore.Repository) Option {
	return func(app *Application) error {
		app.keystoreRepo = repo
		return nil
	}
}

func WithUserRepository(repo user.Repository) Option {
	return func(app *Application) error {
		app.userRepo = repo
		return nil
	}
}

func WithInvitationRepository(repo invitation.Repository) Option {
	return func(app *Application) error {
		app.invitationRepo = repo
		return nil
	}
}

func WithUserService(service user.Service) Option {
	return func(app *Application) error {
		app.userService = service
		return nil
	}
}

func WithKeystoreService(service keystore.Service) Option {
	return func(app *Application) error {
		app.keystoreService = service
		return nil
	}
}

func WithInvitationService(service invitation.Service) Option {
	return func(app *Application) error {
		app.invitationService = service
		return nil
	}
}
