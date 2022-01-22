package application

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/curve25519"

	"myst/pkg/crypto"

	"myst/internal/client/core/remote"

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
	Keystores() ([]*keystore.Keystore, error)
	HealthCheck()
	SignIn(username, password string) error
	SignOut() error
}

type application struct {
	keystoreService keystore.Service
	keystoreRepo    KeystoreRepository
	remote          remote.Client
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

	rc, err := remote.New()
	if err != nil {
		return nil, err
	}

	app.remote = rc

	return app, nil
}

func (app *application) setup() {
	k, err := app.keystoreService.Create(
		keystore.WithName("my-keystore"),
		keystore.WithPassphrase("pass"),
	)
	if err != nil && err.Error() == "already exists" {
		k, err = app.keystoreService.Unlock("0000000000000000000000", "pass")
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < 0; i++ {
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

		err = k.AddEntry(*e1)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = k.AddEntry(*e2)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = k.AddEntry(*e3)
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

	err = app.remote.SignIn("rdnt", "1234")
	if err != nil {
		fmt.Println(err)
		return
	}

	kpath := "data/keystores/" + k.Id() + ".mst"

	b, err := os.ReadFile(kpath)
	if err != nil {
		fmt.Println(err)
		return
	}

	sk, err := app.remote.CreateKeystore(
		"my-keystore", b,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Debug(sk)

	sk2, err := app.remote.Keystore(
		sk.Id(),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Debug(sk2)

	sks, err := app.remote.Keystores()
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Debug(sks)

	u1pub, u1key, err := newKeypair()
	if err != nil {
		fmt.Println(err)
		return
	}

	u2pub, u2key, err := newKeypair()
	if err != nil {
		fmt.Println(err)
		return
	}

	inv, err := app.remote.CreateInvitation(
		"0000000000000000000000", "abcd", u1pub,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	inv, err = app.remote.AcceptInvitation(
		"0000000000000000000000", inv.Id(), u2pub,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	inv, err = app.remote.FinalizeInvitation(
		"0000000000000000000000", inv.Id(), u2pub,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(u1key, u2key)

	log.Debug(inv)
}

func newKeypair() ([]byte, []byte, error) {
	var pub [32]byte
	var key [32]byte

	b, err := crypto.GenerateRandomBytes(32)
	if err != nil {
		return nil, nil, err
	}
	copy(key[:], b)

	curve25519.ScalarBaseMult(&pub, &key)

	return pub[:], key[:], nil
}
