package application

import (
	"myst/pkg/logger"
	"myst/src/client/application/domain/credentials"
	"myst/src/client/application/domain/entry"
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/user"
)

var log = logger.New("app", logger.Blue)

// Enclave is the repository that handles storing and retrieving of the
// user's keystores and credentials. It requires initialization and
// authentication before it can be used. The authentication status can
// expire after some time if the HealthCheck method is not called
// regularly.
type Enclave interface {
	Initialize(password string) error
	IsInitialized() (bool, error)
	Authenticate(password string) error
	HealthCheck()

	CreateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	Keystore(id string) (keystore.Keystore, error)
	UpdateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	Keystores() (map[string]keystore.Keystore, error)
	DeleteKeystore(id string) error

	UpdateCredentials(creds credentials.Credentials) (credentials.Credentials, error)
	Credentials() (credentials.Credentials, error)
}

// Remote is a remote repository that holds upstream enclave/invitations. It is
// used to sync keystores with a remote server in a secure manner, and to
// facilitate inviting users to access keystores or accepting invitations to
// access a keystore from another user. Authenticating with a username and
// password is required to interface with a remote.
type Remote interface {
	Address() string

	CreateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	UpdateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	Keystores(privateKey []byte) (map[string]keystore.Keystore, error)
	DeleteKeystore(id string) error

	CreateInvitation(keystoreRemoteId, inviteeUsername string) (invitation.Invitation, error)
	Invitation(id string) (invitation.Invitation, error)
	AcceptInvitation(id string) (invitation.Invitation, error)
	DeleteInvitation(id string) (invitation.Invitation, error)
	FinalizeInvitation(invitationId string, encryptedKeystoreKey []byte) (invitation.Invitation, error)
	Invitations() (map[string]invitation.Invitation, error)

	Authenticate(username, password string) error
	Register(username, password string, publicKey []byte) (user.User, error)
	Authenticated() bool
	CurrentUser() *user.User
}

type UpdateEntryOptions struct {
	Password *string
	Notes    *string
}

type Application interface {
	Initialize(password string) error
	IsInitialized() (bool, error)
	Authenticate(password string) error
	HealthCheck()

	CreateKeystore(name string) (keystore.Keystore, error)
	DeleteKeystore(id string) error
	Keystore(id string) (keystore.Keystore, error)
	Keystores() (map[string]keystore.Keystore, error)

	CreateEntry(keystoreId string, website, username, password, notes string) (entry.Entry, error)
	UpdateEntry(keystoreId string, entryId string, opts UpdateEntryOptions) (entry.Entry, error)
	DeleteEntry(keystoreId, entryId string) error

	Register(username, password string) (user.User, error)
	CurrentUser() (*user.User, error)

	CreateInvitation(keystoreId string, inviteeUsername string) (invitation.Invitation, error)
	AcceptInvitation(id string) (invitation.Invitation, error)
	DeleteInvitation(id string) (invitation.Invitation, error)
	FinalizeInvitation(invitationId, remoteKeystoreId string, inviteePublicKey []byte) (invitation.Invitation, error)
	Invitations() (map[string]invitation.Invitation, error)
	Invitation(id string) (invitation.Invitation, error)

	Sync() error
}

type application struct {
	enclave Enclave
	remote  Remote
}

func New(opts ...Option) Application {
	app := &application{}

	for _, opt := range opts {
		if opt != nil {
			opt(app)
		}
	}

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
