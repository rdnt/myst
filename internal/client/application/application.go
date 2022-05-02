package application

import (
	"errors"
	"fmt"

	"myst/internal/client/application/domain/entry"
	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/application/keystoreservice"
	"myst/internal/client/remote"
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
	k1, err := app.CreateFirstKeystore("Passwords", "12345678")
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
		_, err = app.CreateKeystoreEntry(k1.Id(), opt...)
		if err != nil {
			panic(err)
		}
	}

	k2, err := app.CreateKeystore("Work")
	if err != nil {
		panic(err)
	}

	_, err = app.CreateKeystoreEntry(k2.Id(),
		entry.WithId("pxnChjAmntT5aG35PM3G12"),
		entry.WithWebsite("www.microsoft.com"), entry.WithUsername("test123@example.com"),
		entry.WithPassword("H278L5qtwvSVs333"),
	)
	if err != nil {
		panic(err)
	}

	_, err = app.CreateKeystore("Other")
	if err != nil {
		panic(err)
	}

	k1, err = app.Keystore(k1.Id())
	if err != nil {
		fmt.Println(err)
		return
	}

	//log.Debug(k1)

	err = app.repositories.remote.SignIn("rdnt", "1234")
	if err != nil {
		fmt.Println(err)
		return
	}

	//jk := enclaverepo.KeystoreToJSON(k1)

	//b, err := json.Marshal(jk)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	key, err := app.KeystoreKey(k1.Id())
	if err != nil {
		fmt.Println(err)
		return
	}

	//log.Debug("KEY", key)

	sk, err := app.repositories.remote.CreateKeystore(
		k1.Name(), key, k1, // TODO: send encrypted keystore with the keystore key (not with the password or the argon2id hash)
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	//log.Debug(sk)

	_, err = app.repositories.remote.Keystore(
		sk.Id,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	//log.Debug(sk2)

	_, err = app.repositories.remote.Keystores()
	if err != nil {
		fmt.Println(err)
		return
	}

	//log.Debug(sks)

	//u1pub, u1key, err := newKeypair()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//u2pub, u2key, err := newKeypair()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//inv, err := app.repositories.remote.CreateInvitation(
	//	"0000000000000000000000", "abcd", u1pub,
	//)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//inv, err = app.repositories.remote.AcceptInvitation(
	//	"0000000000000000000000", inv.Id, u2pub,
	//)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//if inv.InviteeKey == nil {
	//	panic("nil invitee key")
	//}
	//
	//// **** ASYMMETRIC KEY GENERATION ****
	//asymKey, err := curve25519.X25519(u1key, *inv.InviteeKey)
	//if err != nil {
	//	panic(err)
	//}
	//
	//// encrypt the keystore key with the symmetric key
	//b, err := crypto.AES256CBC_Encrypt(asymKey, key)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//inv, err = app.repositories.remote.FinalizeInvitation(
	//	"0000000000000000000000", inv.Id, b,
	//)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//asymKey2, err := curve25519.X25519(u2key, *inv.InviterKey)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//b, err = crypto.AES256CBC_Decrypt(asymKey2, *inv.KeystoreKey)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//if bytes.Compare(b, key) != 0 {
	//	panic("key mismatch")
	//}

	//log.Debug(u1key, u2key)
	//log.Debug(inv)
}

//func newKeypair() ([]byte, []byte, error) {
//	var pub [32]byte
//	var key [32]byte
//
//	b, err := crypto.GenerateRandomBytes(32)
//	if err != nil {
//		return nil, nil, err
//	}
//	copy(key[:], b)
//
//	curve25519.ScalarBaseMult(&pub, &key)
//
//	return pub[:], key[:], nil
//}
