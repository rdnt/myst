package application

import (
	"myst/pkg/logger"
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/keystore/entry"
	"myst/src/client/application/domain/user"
)

var log = logger.New("app", logger.Blue)

type Application interface {
	CreateInvitation(keystoreId string, inviteeUsername string) (invitation.Invitation, error)
	AcceptInvitation(id string) (invitation.Invitation, error)
	DeclineOrCancelInvitation(id string) (invitation.Invitation, error)
	FinalizeInvitation(id string) (invitation.Invitation, error)
	Invitations() (map[string]invitation.Invitation, error)
	Invitation(id string) (invitation.Invitation, error)

	CreateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	DeleteKeystore(id string) error
	Keystore(id string) (keystore.Keystore, error)
	KeystoreByRemoteId(id string) (keystore.Keystore, error)
	CreateKeystoreEntry(keystoreId string, opts ...entry.Option) (entry.Entry, error)
	UpdateKeystoreEntry(keystoreId string, entryId string, password, notes *string) (entry.Entry, error)
	DeleteKeystoreEntry(keystoreId, entryId string) error
	Keystores() (map[string]keystore.Keystore, error)
	UpdateKeystore(k keystore.Keystore) error

	SignIn(username, password string) (user.User, error)
	Register(username, password string) (user.User, error)
	CurrentUser() (*user.User, error)

	HealthCheck()
	Initialize(password string) error
	IsInitialized() (bool, error)
	Authenticate(password string) error

	Sync() error
	Start() error
	Stop() error
	Debug() (map[string]any, error)
}

type application struct {
	enclave Enclave
	remote  Remote
}

func New(opts ...Option) (Application, error) {
	app := &application{}

	for _, opt := range opts {
		err := opt(app)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	return app, nil
}

func (app *application) Start() error {
	defer log.Print("App started")

	return nil
}

func (app *application) Stop() error {
	log.Print("App stopped")

	return nil
}

func (app *application) Debug() (data map[string]any, err error) {
	data = map[string]any{}

	data["keystores"], err = app.enclave.Keystores()
	if err != nil {
		return nil, err
	}

	data["credentials"], err = app.enclave.Credentials()
	if err != nil {
		return nil, err
	}

	data["invitations"], err = app.remote.Invitations()
	if err != nil {
		return nil, err
	}

	return data, nil
}
