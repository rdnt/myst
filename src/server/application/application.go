package application

import (
	"github.com/pkg/errors"

	"myst/pkg/logger"
	"myst/src/server/application/domain/invitation"
	"myst/src/server/application/domain/keystore"
	"myst/src/server/application/domain/user"
)

var log = logger.New("app", logger.Blue)

type InvitationRepository interface {
	CreateInvitation(invitation.Invitation) (invitation.Invitation, error)
	Invitation(id string) (invitation.Invitation, error)
	UpdateInvitation(invitation.Invitation) (invitation.Invitation, error)
	DeleteInvitation(id string) error
	Invitations() ([]invitation.Invitation, error)
}

type KeystoreRepository interface {
	CreateKeystore(keystore.Keystore) (keystore.Keystore, error)
	Keystore(id string) (keystore.Keystore, error)
	UpdateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	DeleteKeystore(id string) error
	Keystores() ([]keystore.Keystore, error)
}

type UserRepository interface {
	CreateUser(user.User) (user.User, error)
	User(id string) (user.User, error)
	UserByUsername(username string) (user.User, error)
	UpdateUser(user.User) (user.User, error)
	Users() ([]user.User, error)
}

var (
	ErrInvitationNotFound = errors.New("invitation not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrNotAllowed         = errors.New("not allowed")
	ErrKeystoreNotFound   = errors.New("keystore not found")
	ErrInvalidUsername    = errors.New("invalid username")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidInvitee     = errors.New("invalid inviter")
)

type KeystoreUpdateParams struct {
	Name    *string
	Payload *[]byte
}

type UserInvitationsOptions struct {
	Status *invitation.Status
}

type Application interface {
	CreateUser(username, password string, publicKey []byte) (user.User, error)
	AuthorizeUser(username, password string) (user.User, error)
	User(id string) (user.User, error)
	UserByUsername(username string) (user.User, error)

	CreateKeystore(name, ownerId string, payload []byte) (keystore.Keystore, error)
	Keystore(id string) (keystore.Keystore, error)
	Keystores() ([]keystore.Keystore, error)
	UpdateKeystore(userId, keystoreId string, params KeystoreUpdateParams) (keystore.Keystore, error)
	DeleteKeystore(userId string, keystoreId string) error
	UserKeystores(userId string) ([]keystore.Keystore, error)

	CreateInvitation(keystoreId, inviterId, inviteeUsername string) (invitation.Invitation, error)
	AcceptInvitation(userId string, invitationId string) (invitation.Invitation, error)
	DeleteInvitation(userId, invitationId string) (invitation.Invitation, error)
	FinalizeInvitation(userId string, invitationId string, encryptedKeystoreKey []byte) (invitation.Invitation, error)
	UserInvitation(userId, invitationId string) (invitation.Invitation, error)
	UserInvitations(userId string, opts UserInvitationsOptions) ([]invitation.Invitation, error)
}

type application struct {
	invitations InvitationRepository
	users       UserRepository
	keystores   KeystoreRepository
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

func WithKeystoreRepository(repo KeystoreRepository) Option {
	return func(app *application) {
		app.keystores = repo
	}
}

func WithUserRepository(repo UserRepository) Option {
	return func(app *application) {
		app.users = repo
	}
}

func WithInvitationRepository(repo InvitationRepository) Option {
	return func(app *application) {
		app.invitations = repo
	}
}
