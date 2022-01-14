package application

import (
	"myst/internal/server/core/domain/invitation"
	"myst/internal/server/core/domain/keystore"
	"myst/internal/server/core/domain/user"
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
