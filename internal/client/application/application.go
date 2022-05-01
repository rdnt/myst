package application

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/curve25519"

	"myst/internal/client/application/domain/entry"
	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/application/keystoreservice"
	"myst/internal/client/remote"
	"myst/pkg/crypto"
	"myst/pkg/logger"
)

var log = logger.New("app", logger.Blue)

var (
	ErrInvalidKeystoreRepository = errors.New("invalid keystore repository")
	ErrInvalidKeystoreService    = errors.New("invalid keystore service")
	ErrAuthenticationFailed      = errors.New("authentiation failed")
)

type Application interface {
	keystore.Service

	Start()
	SignIn(username, password string) error
	SignOut() error
}

type application struct {
	keystore.Service

	repositories struct {
		keystoreRepo keystoreservice.KeystoreRepository
		remote       remote.Client
	}
}

func (app *application) Start() {
	log.Print("App started")

	//app.setup()
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

	if app.repositories.keystoreRepo == nil {
		return nil, ErrInvalidKeystoreRepository
	}

	keystoreService, err := keystoreservice.New(
		keystoreservice.WithKeystoreRepository(app.repositories.keystoreRepo),
	)
	if err != nil {
		panic(err)
	}

	app.Service = keystoreService

	rc, err := remote.New()
	if err != nil {
		return nil, err
	}

	app.repositories.remote = rc

	return app, nil
}

func (app *application) setup() {
	k, err := app.CreateFirstKeystore("my-keystore", "pass")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 0; i++ {
		_, err = app.CreateKeystoreEntry(
			k.Id(),
			entry.WithWebsite("google.com"),
			entry.WithUsername("someuser@google.com"),
			entry.WithPassword("12345678"),
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = app.CreateKeystoreEntry(
			k.Id(),
			entry.WithWebsite("stackoverflow.com"),
			entry.WithUsername("someotheruser@google.com"),
			entry.WithPassword("abcdefghijklmnopqrstuvwxyz"),
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = app.CreateKeystoreEntry(
			k.Id(),
			entry.WithWebsite("reddit.com"),
			entry.WithUsername("somethirduser@yahoo.com"),
			entry.WithPassword("!@*#&$^!@*#&$^!"),
		)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	k, err = app.Keystore(k.Id())
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Debug(k)

	err = app.repositories.remote.SignIn("rdnt", "1234")
	if err != nil {
		fmt.Println(err)
		return
	}

	kpath := "data/keystores/" + k.Id() + ".myst"

	b, err := os.ReadFile(kpath)
	if err != nil {
		fmt.Println(err)
		return
	}

	sk, err := app.repositories.remote.CreateKeystore(
		"my-keystore", b,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Debug(sk)

	sk2, err := app.repositories.remote.Keystore(
		sk.Id(),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Debug(sk2)

	sks, err := app.repositories.remote.Keystores()
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

	inv, err := app.repositories.remote.CreateInvitation(
		"0000000000000000000000", "abcd", u1pub,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	inv, err = app.repositories.remote.AcceptInvitation(
		"0000000000000000000000", inv.Id(), u2pub,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	inv, err = app.repositories.remote.FinalizeInvitation(
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
