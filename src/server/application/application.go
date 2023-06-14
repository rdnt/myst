package application

import (
	"github.com/pkg/errors"

	"myst/src/server/application/domain/invitation"
	"myst/src/server/application/domain/keystore"
	"myst/src/server/application/domain/user"
)

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
	ErrInvitationNotFound   = errors.New("invitation not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrUsernameTaken        = errors.New("username taken")
	ErrForbidden            = errors.New("forbidden")
	ErrKeystoreNotFound     = errors.New("keystore not found")
	ErrInvalidUsername      = errors.New("invalid username")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrAuthenticationFailed = errors.New("authorization failed")
	ErrInvalidInvitee       = errors.New("invalid inviter")
	ErrInviterNotFound      = errors.New("inviter not found")
	ErrInviteeNotFound      = errors.New("invitee not found")
)

type UpdateKeystoreOptions struct {
	Name    *string
	Payload *[]byte
}

type UserInvitationsOptions struct {
	Status *invitation.Status
}

type Application interface {
	CreateUser(username, password string, publicKey []byte) (user.User, error)
	AuthenticateUser(username, password string) (user.User, error)
	User(id string) (user.User, error)
	UserByUsername(username string) (user.User, error)

	CreateKeystore(name, ownerId string, payload []byte) (keystore.Keystore, error)
	Keystore(id string) (keystore.Keystore, error)
	UserKeystore(userId, keystoreId string) (keystore.Keystore, error)
	Keystores() ([]keystore.Keystore, error)
	UpdateKeystore(userId, keystoreId string, opts UpdateKeystoreOptions) (keystore.Keystore, error)
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
