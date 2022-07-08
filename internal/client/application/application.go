package application

import (
	"errors"
	"fmt"
	"time"

	"myst/internal/client/application/domain/enclave"
	"myst/internal/client/application/domain/entry"
	"myst/internal/client/application/domain/invitation"
	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/application/domain/user"
	"myst/internal/client/remote"
	"myst/pkg/logger"
)

var log = logger.New("app", logger.Blue)

var (
	ErrInvalidKeystoreService = errors.New("invalid keystore service")
	ErrAuthenticationFailed   = errors.New("authentiation failed")
)

type Application interface {
	Start() error
	Stop() error

	// remote
	SignIn(username, password string) error
	SignOut() error
	CurrentUser() (*user.User, error)
	Register(username, password string) (user.User, error)
	// CreateKeystore(name string, keystoreKey []byte, keystore *keystore.Keystore) (*generated.Keystore, error)
	// Keystore(id string) (*generated.Keystore, error)
	// Keystores() ([]*generated.Keystore, error)
	// CreateInvitation(keystoreId, inviteeId string) (*generated.Invitation, error)
	// AcceptInvitation(keystoreId, invitationId string) (*generated.Invitation, error)
	// FinalizeInvitation(keystoreId, invitationId string, keystoreKey []byte) (*generated.Invitation, error)

	// keystores
	Authenticate(password string) error
	CreateFirstKeystore(k keystore.Keystore, password string) (keystore.Keystore, error) // TODO: should this be determined during 'CreateKeystore()'?
	CreateEnclave(password string) error
	Enclave() error
	CreateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	DeleteKeystore(id string) error
	Keystore(id string) (keystore.Keystore, error)
	CreateKeystoreEntry(keystoreId string, opts ...entry.Option) (entry.Entry, error)
	UpdateKeystoreEntry(keystoreId string, entryId string, password, notes *string) (entry.Entry, error)
	DeleteKeystoreEntry(keystoreId, entryId string) error
	Keystores() (map[string]keystore.Keystore, error)
	HealthCheck()

	CreateInvitation(keystoreId string, inviteeUsername string) (invitation.Invitation, error)
	AcceptInvitation(id string) (invitation.Invitation, error)
	DeclineOrCancelInvitation(id string) (invitation.Invitation, error)
	FinalizeInvitation(id string) (invitation.Invitation, error)
	Invitations() (map[string]invitation.Invitation, error)
}

type application struct {
	keystores   enclave.KeystoreService
	invitations invitation.Service

	remote remote.Remote
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

	if app.keystores == nil {
		return nil, ErrInvalidKeystoreService
	}

	if app.remote == nil {
		return nil, fmt.Errorf("invalid remote")
	}

	return app, nil
}

func (app *application) Start() error {
	defer log.Print("App started")

	go func() {
		for {
			if app.remote.SignedIn() {
				err := app.sync()
				if err != nil {
					log.Error(err)
				}
				// err := app.remote.SignIn()
				// if err != nil {
				//	log.Error(err)
				// }
				//
				// if app.remote.SignedIn() {
				//	err := app.sync()
				//	if err != nil {
				//		log.Error(err)
				//	}
				// }
			}

			time.Sleep(10 * time.Second)
		}
	}()

	return nil
}

func (app *application) Stop() error {
	log.Print("App stopped")

	return nil
}

func (app *application) HealthCheck() {
	app.keystores.HealthCheck()
}

func (app *application) setup() {
	k1, err := app.keystores.CreateFirstKeystore(keystore.New(keystore.WithName("Passwords")), "12345678")
	if err != nil {
		panic(err)
	}

	opts := [][]entry.Option{
		{
			entry.WithId("px5VAUMgPMBtjrAj9ajeFR"),
			entry.WithWebsite("github.com"), entry.WithUsername("rdntdev@gmail.com"),
			entry.WithPassword("nzK&d#u+MjFU8p&4UhL)s3+h"),
			entry.WithNotes("Lorem ipsum"),
		},
		{
			entry.WithId("Vxg4iMtmXUw76t77hb6m3B"),
			entry.WithWebsite("youtube.com"), entry.WithUsername("oldsnut@gmailni.com"),
			entry.WithPassword("tsksWgABXhvh9LfF"),
			entry.WithNotes("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."),
		},
		{
			entry.WithId("xKPddK9gtAbUT3Ej93ShZ"),
			entry.WithWebsite("facebook.com"), entry.WithUsername("pete24uk@test130.com"),
			entry.WithPassword("uXekxDRk6bmvvpda"),
			entry.WithNotes("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."),
		},
		{
			entry.WithId("ZByasmj3aLHgRMJDeDiXS4"),
			entry.WithWebsite("baidu.com"), entry.WithUsername("swissly0@telemol.online"),
			entry.WithPassword("6Yu4k2YkNPxpZkHn"),
		},
		{
			entry.WithId("hozhDrCZGqjZ2VGLcpcuNi"),
			entry.WithWebsite("yahoo.com"), entry.WithUsername("swissly0@telemol.online"),
			entry.WithPassword("LxRdhTTSg4Adfkc4"),
		},
		{
			entry.WithId("tEsfGypbPAVCNMAKWtw2mD"),
			entry.WithWebsite("amazon.com"), entry.WithUsername("manhosobpe15@zetgets.com"),
			entry.WithPassword("qHcZsZxPf8acHxxA"),
		},
		{
			entry.WithId("McB6akkM3C5XpzXfMYhasU"),
			entry.WithWebsite("wikipedia.org"), entry.WithUsername("chuninoleg1971@piftir.com"),
			entry.WithPassword("DT3sftJuRjxWFg68"),
		},
		{
			entry.WithId("YBS32eK8XbeV6ujaY5xERK"),
			entry.WithWebsite("twitter.com"), entry.WithUsername("ninablackangel@test.com"),
			entry.WithPassword("ndUZ6KGduD53up4R"),
		},
		{
			entry.WithId("Fy7HDsbQqkYsbevjuqSG65"),
			entry.WithWebsite("bbc.com"), entry.WithUsername("kgdlove@omdlism.com"),
			entry.WithPassword("jy9EpWExSmmtHa6g"),
		},
		{
			entry.WithId("r5TbidUGZkZeqbP7iCySBn"),
			entry.WithWebsite("steampowered.com"), entry.WithUsername("totinoprato@roselarose.com"),
			entry.WithPassword("tbRCJ9uHvxLm9S5q"),
		},
		{
			entry.WithId("pxnChjAmntT5aG35PM3GL4"),
			entry.WithWebsite("bing.com"), entry.WithUsername("tbiggs@massageshophome.com"),
			entry.WithPassword("H278L5qtwvSVsQzt"),
		},
	}

	for _, opt := range opts {
		_, err = app.keystores.CreateKeystoreEntry(k1.Id, opt...)
		if err != nil {
			panic(err)
		}
	}

	k2, err := app.keystores.CreateKeystore(keystore.New(keystore.WithName("Work")))
	if err != nil {
		panic(err)
	}

	_, err = app.keystores.CreateKeystoreEntry(k2.Id,
		entry.WithId("pxnChjAmntT5aG35PM3G12"),
		entry.WithWebsite("www.microsoft.com"), entry.WithUsername("test123@example.com"),
		entry.WithPassword("H278L5qtwvSVs333"),
	)
	if err != nil {
		panic(err)
	}

	_, err = app.keystores.CreateKeystore(keystore.New(keystore.WithName("Other")))
	if err != nil {
		panic(err)
	}

	k1, err = app.keystores.Keystore(k1.Id)
	if err != nil {
		fmt.Println(err)
		return
	}

	// err = app.remote.SignIn()
	// if err != nil {
	//	fmt.Println(err)
	//	return
	// }

	// key, err := app.keystores.EncryptedKeystoreKey(k1.Id())
	// if err != nil {
	//	fmt.Println(err)
	//	return
	// }

	// TODO: fix applicationrefactor debug initialization
	// sk, err := app.remote.UploadKeystore(
	//	k1.Name(), key, k1, // TODO: send encrypted keystore with the keystore key (not with the password or the argon2id hash)
	// )
	// if err != nil {
	//	fmt.Println(err)
	//	return
	// }
	//
	// _, err = app.remote.Keystore(
	//	sk.Id,
	// )
	// if err != nil {
	//	fmt.Println(err)
	//	return
	// }
	//
	// _, err = app.remote.Keystores()
	// if err != nil {
	//	fmt.Println(err)
	//	return
	// }
}
