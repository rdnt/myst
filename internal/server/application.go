package application

import (
	"errors"

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
	userService       user.Service
	keystoreService   keystore.Service
	invitationService invitation.Service
	invitationRepo    invitation.Repository
	userRepo          user.Repository
	keystoreRepo      keystore.Repository
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

	if app.keystoreRepo == nil {
		return nil, ErrInvalidKeystoreRepository
	}

	if app.userRepo == nil {
		return nil, ErrInvalidUserRepository
	}

	if app.invitationRepo == nil {
		return nil, ErrInvalidUserRepository
	}

	if app.userService == nil {
		return nil, ErrInvalidUserService
	}

	if app.keystoreService == nil {
		return nil, ErrInvalidUserService
	}

	if app.invitationService == nil {
		return nil, ErrInvalidUserService
	}

	return app, nil
}

func (app *Application) setup() {
	u, err := app.userService.Register(
		user.WithUsername("rdnt"),
		user.WithPassword("1234"),
	)
	if err != nil {
		panic(err)
	}

	log.Debug(u)

	u2, err := app.userService.Register(
		user.WithUsername("abcd"),
		user.WithPassword("5678"),
	)
	if err != nil {
		panic(err)
	}

	log.Debug(u2)

	k, err := app.keystoreService.Create(
		keystore.WithName("my-keystore"),
		keystore.WithKeystore([]byte("payload")),
		keystore.WithOwner(u),
	)

	log.Debug(k)

	inv, err := app.invitationService.Create(
		invitation.WithInviter(u),
		invitation.WithKeystore(k),
		invitation.WithInvitee(u2),
		invitation.WithInviterKey([]byte("inviter-key")),
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
