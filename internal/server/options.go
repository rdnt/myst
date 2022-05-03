package application

import (
	"myst/internal/server/core/domain/invitation"
	"myst/internal/server/core/domain/keystore"
	"myst/internal/server/core/domain/user"
)

type Option func(app *application) error

func WithKeystoreRepository(repo keystore.Repository) Option {
	return func(app *application) error {
		app.repositories.keystoreRepo = repo
		return nil
	}
}

func WithUserRepository(repo user.Repository) Option {
	return func(app *application) error {
		app.repositories.userRepo = repo
		return nil
	}
}

func WithInvitationRepository(repo invitation.Repository) Option {
	return func(app *application) error {
		app.repositories.invitationRepo = repo
		return nil
	}
}
