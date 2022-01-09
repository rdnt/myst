package application

import (
	"errors"
	"fmt"

	"myst/internal/client/core/domain/keystore"
	"myst/internal/client/core/domain/keystore/entry"

	"myst/pkg/logger"
)

var log = logger.New("app", logger.Blue)

var (
	ErrInvalidKeystoreRepository = errors.New("invalid keystore repository")
	ErrInvalidKeystoreService    = errors.New("invalid keystore service")
)

type Application interface {
	Start()
	CreateKeystore(name string, passphrase string) (*keystore.Keystore, error)
	UnlockKeystore(keystoreId string, passphrase string) (*keystore.Keystore, error)
	UpdateKeystore(k *keystore.Keystore) error
	Keystore(id string) (*keystore.Keystore, error)
	HealthCheck()
	SignIn(username, password string) error
	SignOut() error
}

type application struct {
	keystoreService keystore.Service
	keystoreRepo    KeystoreRepository
	//serverHttpClient TODO
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
		return nil, ErrInvalidKeystoreService
	}

	return app, nil
}

func (app *application) setup() {
	k, err := app.keystoreService.Create(
		keystore.WithName("my-keystore"),
		keystore.WithPassphrase("pass"),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 3; i++ {
		e1, err := entry.New(
			entry.WithLabel("google.com"),
			entry.WithUsername("someuser@google.com"),
			entry.WithPassword("12345678"),
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		e2, err := entry.New(
			entry.WithLabel("stackoverflow.com"),
			entry.WithUsername("someotheruser@google.com"),
			entry.WithPassword("abcdefghijklmnopqrstuvwxyz"),
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		e3, err := entry.New(
			entry.WithLabel("reddit.com"),
			entry.WithUsername("somethirduser@yahoo.com"),
			entry.WithPassword("!@*#&$^!@*#&$^!"),
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = k.AddEntry(e1)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = k.AddEntry(e2)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = k.AddEntry(e3)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	err = app.keystoreService.Update(k)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Debug(k)
}
