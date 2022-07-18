package application

import (
	"time"

	"myst/pkg/logger"
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/remote"
	"myst/src/client/application/domain/user"
)

var log = logger.New("app", logger.Blue)

type Application interface {
	invitation.Service
	keystore.Service
	user.Service
	Repository

	// Authenticate(password string) error
	// HealthCheck()
	// Initialize(password string) error
	// IsInitialized() error

	Start() error
	Stop() error
	Debug() (map[string]any, error)
}

type application struct {
	keystores   keystore.Repository
	invitations invitation.Repository
	credentials remote.Repository
	repo        Repository

	remote Remote
}

func New(opts ...Option) (Application, error) {
	app := &application{}

	for _, opt := range opts {
		err := opt(app)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
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

			time.Sleep(5 * time.Second)
		}
	}()

	return nil
}

func (app *application) Stop() error {
	log.Print("App stopped")

	return nil
}

func (app *application) Debug() (data map[string]any, err error) {
	data = map[string]any{}

	data["keystores"], err = app.keystores.Keystores()
	if err != nil {
		return nil, err
	}

	data["credentials"], err = app.credentials.Remote()
	if err != nil {
		return nil, err
	}

	data["invitations"], err = app.invitations.Invitations()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// func (app *application) setup() {
// 	k1, err := app.CreateFirstKeystore(keystore.New(keystore.WithName("Passwords")), "12345678")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	opts := [][]entry.Option{
// 		{
// 			entry.WithId("px5VAUMgPMBtjrAj9ajeFR"),
// 			entry.WithWebsite("github.com"), entry.WithUsername("rdntdev@gmail.com"),
// 			entry.WithPassword("nzK&d#u+MjFU8p&4UhL)s3+h"),
// 			entry.WithNotes("Lorem ipsum"),
// 		},
// 		{
// 			entry.WithId("Vxg4iMtmXUw76t77hb6m3B"),
// 			entry.WithWebsite("youtube.com"), entry.WithUsername("oldsnut@gmailni.com"),
// 			entry.WithPassword("tsksWgABXhvh9LfF"),
// 			entry.WithNotes("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."),
// 		},
// 		{
// 			entry.WithId("xKPddK9gtAbUT3Ej93ShZ"),
// 			entry.WithWebsite("facebook.com"), entry.WithUsername("pete24uk@test130.com"),
// 			entry.WithPassword("uXekxDRk6bmvvpda"),
// 			entry.WithNotes("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."),
// 		},
// 		{
// 			entry.WithId("ZByasmj3aLHgRMJDeDiXS4"),
// 			entry.WithWebsite("baidu.com"), entry.WithUsername("swissly0@telemol.online"),
// 			entry.WithPassword("6Yu4k2YkNPxpZkHn"),
// 		},
// 		{
// 			entry.WithId("hozhDrCZGqjZ2VGLcpcuNi"),
// 			entry.WithWebsite("yahoo.com"), entry.WithUsername("swissly0@telemol.online"),
// 			entry.WithPassword("LxRdhTTSg4Adfkc4"),
// 		},
// 		{
// 			entry.WithId("tEsfGypbPAVCNMAKWtw2mD"),
// 			entry.WithWebsite("amazon.com"), entry.WithUsername("manhosobpe15@zetgets.com"),
// 			entry.WithPassword("qHcZsZxPf8acHxxA"),
// 		},
// 		{
// 			entry.WithId("McB6akkM3C5XpzXfMYhasU"),
// 			entry.WithWebsite("wikipedia.org"), entry.WithUsername("chuninoleg1971@piftir.com"),
// 			entry.WithPassword("DT3sftJuRjxWFg68"),
// 		},
// 		{
// 			entry.WithId("YBS32eK8XbeV6ujaY5xERK"),
// 			entry.WithWebsite("twitter.com"), entry.WithUsername("ninablackangel@test.com"),
// 			entry.WithPassword("ndUZ6KGduD53up4R"),
// 		},
// 		{
// 			entry.WithId("Fy7HDsbQqkYsbevjuqSG65"),
// 			entry.WithWebsite("bbc.com"), entry.WithUsername("kgdlove@omdlism.com"),
// 			entry.WithPassword("jy9EpWExSmmtHa6g"),
// 		},
// 		{
// 			entry.WithId("r5TbidUGZkZeqbP7iCySBn"),
// 			entry.WithWebsite("steampowered.com"), entry.WithUsername("totinoprato@roselarose.com"),
// 			entry.WithPassword("tbRCJ9uHvxLm9S5q"),
// 		},
// 		{
// 			entry.WithId("pxnChjAmntT5aG35PM3GL4"),
// 			entry.WithWebsite("bing.com"), entry.WithUsername("tbiggs@massageshophome.com"),
// 			entry.WithPassword("H278L5qtwvSVsQzt"),
// 		},
// 	}
//
// 	for _, opt := range opts {
// 		_, err = app.CreateKeystoreEntry(k1.Id, opt...)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
//
// 	k2, err := app.CreateKeystore(keystore.New(keystore.WithName("Work")))
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	_, err = app.CreateKeystoreEntry(k2.Id,
// 		entry.WithId("pxnChjAmntT5aG35PM3G12"),
// 		entry.WithWebsite("www.microsoft.com"), entry.WithUsername("test123@example.com"),
// 		entry.WithPassword("H278L5qtwvSVs333"),
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	_, err = app.CreateKeystore(keystore.New(keystore.WithName("Other")))
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	k1, err = app.Keystore(k1.Id)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
//
// 	// err = app.remote.SignIn()
// 	// if err != nil {
// 	//	fmt.Println(err)
// 	//	return
// 	// }
//
// 	// key, err := app.keystores.EncryptedKeystoreKey(k1.Id())
// 	// if err != nil {
// 	//	fmt.Println(err)
// 	//	return
// 	// }
//
// 	// TODO: fix applicationrefactor debug initialization
// 	// sk, err := app.remote.UploadKeystore(
// 	//	k1.Name(), key, k1, // TODO: send encrypted keystore with the keystore key (not with the password or the argon2id hash)
// 	// )
// 	// if err != nil {
// 	//	fmt.Println(err)
// 	//	return
// 	// }
// 	//
// 	// _, err = app.remote.Keystore(
// 	//	sk.Id,
// 	// )
// 	// if err != nil {
// 	//	fmt.Println(err)
// 	//	return
// 	// }
// 	//
// 	// _, err = app.remote.Keystores()
// 	// if err != nil {
// 	//	fmt.Println(err)
// 	//	return
// 	// }
// }
