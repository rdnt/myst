package application

import (
	"errors"

	"myst/internal/server/core/invitationservice"
	"myst/internal/server/core/keystoreservice"
	"myst/internal/server/core/userservice"

	"myst/internal/server/core/domain/invitation"
	"myst/internal/server/core/domain/keystore"
	"myst/internal/server/core/domain/user"

	"myst/pkg/logger"
)

var log = logger.New("app", logger.Blue)

var (
	ErrInvalidKeystoreRepository = errors.New("invalid keystore repository")
	ErrInvalidUserRepository     = errors.New("invalid user repository")
	ErrInvalidUserService        = errors.New("invalid user service")
)

type Application struct {
	invitationRepo invitation.Repository
	userRepo       user.Repository
	keystoreRepo   keystore.Repository

	Users       user.Service
	Keystores   keystore.Service
	Invitations invitation.Service
}

func (app *Application) Start() {
	log.Print("App started")

	app.setup()
}

func New(opts ...Option) (*Application, error) {
	app := &Application{}

	for _, opt := range opts {
		err := opt(app)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	var err error

	app.Users, err = userservice.New(
		userservice.WithUserRepository(app.userRepo),
		userservice.WithKeystoreRepository(app.keystoreRepo),
	)
	if err != nil {
		return nil, err
	}

	app.Keystores, err = keystoreservice.New(
		keystoreservice.WithUserRepository(app.userRepo),
		keystoreservice.WithKeystoreRepository(app.keystoreRepo),
	)
	if err != nil {
		return nil, err
	}

	app.Invitations, err = invitationservice.New(
		invitationservice.WithUserRepository(app.userRepo),
		invitationservice.WithKeystoreRepository(app.keystoreRepo),
		invitationservice.WithInvitationRepository(app.invitationRepo),
	)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *Application) setup() {
	u, err := app.Users.Register(
		user.WithUsername("rdnt"),
		user.WithPassword("1234"),
	)
	if err != nil {
		panic(err)
	}

	log.Debug(u)

	u2, err := app.Users.Register(
		user.WithUsername("abcd"),
		user.WithPassword("5678"),
	)
	if err != nil {
		panic(err)
	}

	log.Debug(u2)

	k, err := app.Keystores.Create(
		keystore.WithName("my-keystore"),
		keystore.WithKeystore([]byte("payload")),
		keystore.WithOwner(*u),
	)

	log.Debug(k)

	inv, err := app.Invitations.Create(
		k.Id(), u.Id(), u2.Id(), []byte("inviter-key"),
	)
	if err != nil {
		panic(err)
	}

	log.Debug("created invitation")
	log.Debug(inv)

	err = inv.Accept([]byte("invitee-key"))
	if err != nil {
		panic(err)
	}

	err = app.invitationRepo.Update(inv)
	if err != nil {
		panic(err)
	}

	log.Debug("accepted invitation")
	log.Debug(inv)

	err = inv.Finalize([]byte("keystore-key"))
	if err != nil {
		panic(err)
	}

	err = app.invitationRepo.Update(inv)
	if err != nil {
		panic(err)
	}

	log.Debug("finalized invitation")
	log.Debug(inv)
}
