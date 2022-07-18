package application

import (
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/remote"
)

type Option func(app *application) error

func WithRemote(remote Remote) Option {
	return func(app *application) error {
		app.remote = remote
		return nil
	}
}

func WithKeystoreRepository(repo keystore.Repository) Option {
	return func(app *application) error {
		app.keystores = repo
		return nil
	}
}

func WithInvitationRepository(repo invitation.Repository) Option {
	return func(app *application) error {
		app.invitations = repo
		return nil
	}
}

func WithRepository(repo Repository) Option {
	return func(app *application) error {
		app.repo = repo
		return nil
	}
}

func WithCredentials(creds remote.Repository) Option {
	return func(app *application) error {
		app.credentials = creds
		return nil
	}
}
