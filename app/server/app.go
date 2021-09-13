package server

import (
	"errors"
	"fmt"

	"myst/app/server/domain/keystore"
	"myst/app/server/domain/user"
	"myst/pkg/logger"
)

var (
	ErrInvalidKeystoreRepository = errors.New("invalid keystore repository")
	ErrInvalidUserRepository     = errors.New("invalid user repository")
	ErrInvalidUserService        = errors.New("invalid user service")
)

type Application struct {
	keystorerepo    keystore.Repository
	userrepo        user.Repository
	userService     user.Service
	keystoreService keystore.Service
}

func (app *Application) Start() {
	u, err := app.userService.Register(
		user.WithUsername("rdnt"),
		user.WithPassword("1234"),
	)
	if err != nil {
		panic(err)
	}

	k, err := app.keystoreService.Create(
		keystore.WithName("my-keystore"),
		keystore.WithKeystore([]byte("payload")),
		keystore.WithOwner(u),aewe
	)

	fmt.Println(u)
	fmt.Println(k)
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

	if app.keystorerepo == nil {
		return nil, ErrInvalidKeystoreRepository
	}

	if app.userrepo == nil {
		return nil, ErrInvalidUserRepository
	}

	if app.userService == nil {
		return nil, ErrInvalidUserService
	}

	return app, nil
}
