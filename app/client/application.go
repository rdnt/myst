package client

import (
	"errors"

	"myst/app/client/core/domain/keystore"
	"myst/pkg/logger"
)

var log = logger.New("app", logger.Blue)

var (
	ErrInvalidKeystoreRepository = errors.New("invalid keystore repository")
	ErrInvalidUserRepository     = errors.New("invalid user repository")
	ErrInvalidUserService        = errors.New("invalid user service")
)

type Application interface {
	Start()
	CreateKeystore(name string, passphrase string) (*keystore.Keystore, error)
	UnlockKeystore(keystoreId string, passphrase string) (*keystore.Keystore, error)
	Keystore(id string) (*keystore.Keystore, error)
	HealthCheck()
}

type application struct {
	keystoreService keystore.Service
	keystoreRepo    KeystoreRepository
}

type KeystoreRepository interface {
	keystore.Repository
	Unlock(keystoreId string, passphrase string) (*keystore.Keystore, error)
	HealthCheck()
}

func (app *application) Start() {
	log.Print("App started")

	app.setup()
}

func New(opts ...Option) (*application, error) {
	app := &application{}

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

	if app.keystoreService == nil {
		return nil, ErrInvalidUserService
	}

	return app, nil
}

func (app *application) setup() {
	k, err := app.keystoreService.Create(
		keystore.WithName("my-keystore"),
		keystore.WithPassphrase("my-passphrase"),
	)

	if err != nil {
		panic(err)
	}

	log.Debug(k)
}
