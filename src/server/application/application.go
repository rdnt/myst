package application

import (
	"myst/pkg/logger"
	"myst/src/server/application/domain/invitation"
	"myst/src/server/application/domain/keystore"
	"myst/src/server/application/domain/user"
)

var log = logger.New("app", logger.Blue)

type Application interface {
	user.Service
	keystore.Service
	invitation.Service

	Start() error
	Stop() error
	Debug() (map[string]any, error)
	DebugUpdateUserPublicKey(userId string, publicKey []byte) error
}

type application struct {
	invitations invitation.Repository
	users       user.Repository
	keystores   keystore.Repository
}

func New(opts ...Option) (Application, error) {
	app := &application{}

	for _, opt := range opts {
		opt(app)
	}

	return app, nil
}

type Option func(app *application)

func WithKeystoreRepository(repo keystore.Repository) Option {
	return func(app *application) {
		app.keystores = repo
	}
}

func WithUserRepository(repo user.Repository) Option {
	return func(app *application) {
		app.users = repo
	}
}

func WithInvitationRepository(repo invitation.Repository) Option {
	return func(app *application) {
		app.invitations = repo
	}
}

func (app *application) Start() error {
	log.Print("App started")

	// app.setup()

	return nil
}

// func (app *application) setup() {
// 	u, err := app.CreateUser("rdnt", "1234", []byte("rdntPublicKey"))
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	log.Debug(u)
//
// 	u, err = app.CreateUser("abcd", "5678", []byte("abcdPublicKey"))
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	log.Debug(u)
// }

func (app *application) Stop() error {
	log.Print("App stopped")

	return nil
}

func (app *application) Debug() (data map[string]any, err error) {
	data = map[string]any{}

	data["users"], err = app.users.Users()
	if err != nil {
		return nil, err
	}

	data["keystores"], err = app.keystores.Keystores()
	if err != nil {
		return nil, err
	}

	data["invitations"], err = app.invitations.Invitations()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (app *application) DebugUpdateUserPublicKey(userId string, publicKey []byte) error {
	u, err := app.users.User(userId)
	if err != nil {
		return err
	}

	u.PublicKey = publicKey

	_, err = app.users.UpdateUser(u)
	if err != nil {
		return err
	}

	return nil
}
