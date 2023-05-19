package application

import (
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
	Invitations() ([]invitation.Invitation, error)
}

type KeystoreRepository interface {
	CreateKeystore(keystore.Keystore) (keystore.Keystore, error)
	Keystore(id string) (keystore.Keystore, error)
	UpdateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	Keystores() ([]keystore.Keystore, error)
}

type UserRepository interface {
	CreateUser(user.User) (user.User, error)
	User(id string) (user.User, error)
	UserByUsername(username string) (user.User, error)
	UpdateUser(user.User) (user.User, error)
	Users() ([]user.User, error)
}

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
	UserKeystore(userId, keystoreId string) (keystore.Keystore, error)
	UserKeystores(userId string) ([]keystore.Keystore, error)

	CreateInvitation(keystoreId, inviterId, inviteeUsername string) (invitation.Invitation, error)
	AcceptInvitation(userId string, invitationId string) (invitation.Invitation, error)
	DeclineOrCancelInvitation(userId, invitationId string) (invitation.Invitation, error)
	FinalizeInvitation(invitationId string, encryptedKeystoreKey []byte) (invitation.Invitation, error)
	UserInvitation(userId, invitationId string) (invitation.Invitation, error)
	UserInvitations(userId string, opts UserInvitationsOptions) ([]invitation.Invitation, error)

	Start() error
	Stop() error
	Debug() (map[string]any, error)
	DebugUpdateUserPublicKey(userId string, publicKey []byte) error
}

type application struct {
	invitations InvitationRepository
	users       UserRepository
	keystores   KeystoreRepository
}

func New(opts ...Option) (Application, error) {
	app := &application{}

	for _, opt := range opts {
		opt(app)
	}

	return app, nil
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

	// data["users"], err = app.users.Users()
	// if err != nil {
	// 	return nil, err
	// }
	//
	// data["keystores"], err = app.keystores.Keystores()
	// if err != nil {
	// 	return nil, err
	// }
	//
	// data["invitations"], err = app.invitations.Invitations()
	// if err != nil {
	// 	return nil, err
	// }

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
