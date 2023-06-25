package application

import (
	"sync"
	"time"

	"myst/pkg/logger"
	"myst/src/client/application/domain/entry"
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/user"
)

var log = logger.New("app", logger.Blue)

type UpdateEntryOptions struct {
	Password *string
	Notes    *string
}

type Application interface {
	Initialize(password string) (sessionId []byte, err error)
	IsInitialized(sessionId []byte) (bool, error)
	Authenticate(password string) (sessionId []byte, err error)
	HealthCheck(sessionId []byte) error

	CreateKeystore(sessionId []byte, name string) (keystore.Keystore, error)
	DeleteKeystore(sessionId []byte, id string) error
	Keystore(sessionId []byte, id string) (keystore.Keystore, error)
	Keystores(sessionId []byte) (map[string]keystore.Keystore, error)

	CreateEntry(sessionId []byte, keystoreId string, website, username, password, notes string) (entry.Entry, error)
	UpdateEntry(sessionId []byte, keystoreId string, entryId string, opts UpdateEntryOptions) (entry.Entry, error)
	DeleteEntry(sessionId []byte, keystoreId, entryId string) error

	Register(sessionId []byte, username, password string) (user.User, error)
	CurrentUser(sessionId []byte) (*user.User, error)
	SharedSecret(sessionId []byte, userId string) ([]byte, error)

	CreateInvitation(sessionId []byte, keystoreId string, inviteeUsername string) (invitation.Invitation, error)
	AcceptInvitation(sessionId []byte, id string) (invitation.Invitation, error)
	DeleteInvitation(sessionId []byte, id string) (invitation.Invitation, error)
	FinalizeInvitation(sessionId []byte, invitationId, remoteKeystoreId string, inviteePublicKey []byte) (invitation.Invitation, error)
	Invitations(sessionId []byte) (map[string]invitation.Invitation, error)
	Invitation(sessionId []byte, id string) (invitation.Invitation, error)

	Sync() error
}

type application struct {
	enclave  Enclave
	remote   Remote
	sessions map[string]time.Time
	key      []byte
	mux      sync.Mutex
}

func New(opts ...Option) Application {
	app := &application{
		sessions: make(map[string]time.Time),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(app)
		}
	}

	go app.startHealthCheck()

	return app
}

type Option func(app *application)

func WithRemote(remote Remote) Option {
	return func(app *application) {
		app.remote = remote
	}
}

func WithEnclave(enclave Enclave) Option {
	return func(app *application) {
		app.enclave = enclave
	}
}
