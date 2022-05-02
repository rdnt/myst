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
	repositories struct {
		invitationRepo invitation.Repository
		userRepo       user.Repository
		keystoreRepo   keystore.Repository
	}

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
		userservice.WithUserRepository(app.repositories.userRepo),
		userservice.WithKeystoreRepository(app.repositories.keystoreRepo),
	)
	if err != nil {
		return nil, err
	}

	app.Keystores, err = keystoreservice.New(
		keystoreservice.WithUserRepository(app.repositories.userRepo),
		keystoreservice.WithKeystoreRepository(app.repositories.keystoreRepo),
	)
	if err != nil {
		return nil, err
	}

	app.Invitations, err = invitationservice.New(
		invitationservice.WithUserRepository(app.repositories.userRepo),
		invitationservice.WithKeystoreRepository(app.repositories.keystoreRepo),
		invitationservice.WithInvitationRepository(app.repositories.invitationRepo),
	)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *Application) setup() {
	_, err := app.Users.Register(
		user.WithUsername("rdnt"),
		user.WithPassword("1234"),
	)
	if err != nil {
		panic(err)
	}

	//log.Debug(u)

	_, err = app.Users.Register(
		user.WithUsername("abcd"),
		user.WithPassword("5678"),
	)
	if err != nil {
		panic(err)
	}

	//log.Debug(u2)

	//k, err := app.Keystores.Create("my-keystore", u.Id(), []byte("payload"))
	//if err != nil {
	//	panic(err)
	//}

	//log.Debug(k)
	//
	//u.OwnKeystore(k.Id())

	//err = app.userRepo.Update(u)
	//if err != nil {
	//	panic(err)
	//}
	//
	//log.Debug(u)
	//
	//inv, err := app.Invitations.Create(
	//	k.Id(), u.Id(), u2.Id(), []byte("inviter-key"),
	//)
	//if err != nil {
	//	panic(err)
	//}
	//
	//log.Debug("created invitation")
	//log.Debug(inv)
	//
	//err = inv.Accept([]byte("invitee-key"))
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = app.invitationRepo.Update(inv)
	//if err != nil {
	//	panic(err)
	//}
	//
	//log.Debug("accepted invitation")
	//log.Debug(inv)
	//
	//err = inv.Finalize([]byte("keystore-key"))
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = app.invitationRepo.Update(inv)
	//if err != nil {
	//	panic(err)
	//}
	//
	//log.Debug("finalized invitation")
	//log.Debug(inv)
}
