package keystoreservice

import (
	"errors"
	"myst/internal/client/application/domain/entry"
	"myst/internal/client/application/domain/keystore"
	"myst/pkg/logger"
	"strings"
)

var (
	ErrInvalidKeystoreRepository = errors.New("invalid keystore repository")
	ErrAuthenticationRequired    = errors.New("authentication required")
	ErrAuthenticationFailed      = errors.New("authentication failed")
	ErrInitializationRequired    = errors.New("initialization required")
	ErrInvalidPassword           = errors.New("invalid password")
	ErrEntryNotFound             = errors.New("entry not found")
)

type KeystoreRepository interface {
	keystore.Repository
	Authenticate(password string) error
	Initialize(password string) error
	HealthCheck()
}

type service struct {
	keystoreRepo KeystoreRepository
}

func (s *service) KeystoreEntries(id string) (map[string]entry.Entry, error) {
	k, err := s.keystoreRepo.Keystore(id)
	if err != nil {
		return nil, err
	}

	return k.Entries(), nil
}

func (s *service) UpdateKeystoreEntry(keystoreId, entryId string, password, notes *string) (entry.Entry, error) {
	// do not allow empty password
	if password != nil && strings.TrimSpace(*password) == "" {
		return entry.Entry{}, ErrInvalidPassword
	}

	k, err := s.keystoreRepo.Keystore(keystoreId)
	if err != nil {
		return entry.Entry{}, err
	}

	entries := k.Entries()

	e, ok := entries[entryId]
	if !ok {
		return entry.Entry{}, ErrEntryNotFound
	}

	if password != nil {
		e.SetPassword(*password)
	}

	if notes != nil {
		e.SetNotes(*notes)
	}

	entries[e.Id()] = e

	k.SetEntries(entries)

	return e, s.UpdateKeystore(k)
}

func (s *service) CreateKeystore(name string) (*keystore.Keystore, error) {
	return s.keystoreRepo.Create(keystore.WithName(name))
}

func (s *service) CreateKeystoreEntry(keystoreId string, opts ...entry.Option) (entry.Entry, error) {
	k, err := s.keystoreRepo.Keystore(keystoreId)
	if err != nil {
		return entry.Entry{}, err
	}

	e := entry.New(opts...)

	entries := k.Entries()
	entries[e.Id()] = e
	k.SetEntries(entries)

	return e, s.keystoreRepo.Update(k)
}

func (s *service) CreateFirstKeystore(name, password string) (*keystore.Keystore, error) {
	err := s.keystoreRepo.Initialize(password)
	if err != nil {
		return nil, err
	}

	// TODO: remove dummy keystores and properly return error
	s.keystoreRepo.Create(
		keystore.WithName(name), keystore.WithEntries(
			map[string]entry.Entry{
				"px5VAUMgPMBtjrAj9ajeFR": entry.New(
					entry.WithId("px5VAUMgPMBtjrAj9ajeFR"),
					entry.WithWebsite("github.com"), entry.WithUsername("rdntdev@gmail.com"),
					entry.WithPassword("nzK&d#u+MjFU8p&4UhL)s3+h"),
					entry.WithNotes("Lorem ipsum"),
				),
				"Vxg4iMtmXUw76t77hb6m3B": entry.New(
					entry.WithId("Vxg4iMtmXUw76t77hb6m3B"),
					entry.WithWebsite("youtube.com"), entry.WithUsername("oldsnut@gmailni.com"),
					entry.WithPassword("tsksWgABXhvh9LfF"),
					entry.WithNotes("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."),
				),
				"xKPddK9gtAbUT3Ej93ShZ": entry.New(
					entry.WithId("xKPddK9gtAbUT3Ej93ShZ"),
					entry.WithWebsite("facebook.com"), entry.WithUsername("pete24uk@test130.com"),
					entry.WithPassword("uXekxDRk6bmvvpda"),
					entry.WithNotes("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."),
				),
				"ZByasmj3aLHgRMJDeDiXS4": entry.New(
					entry.WithId("ZByasmj3aLHgRMJDeDiXS4"),
					entry.WithWebsite("baidu.com"), entry.WithUsername("swissly0@telemol.online"),
					entry.WithPassword("6Yu4k2YkNPxpZkHn"),
				),
				"hozhDrCZGqjZ2VGLcpcuNi": entry.New(
					entry.WithId("hozhDrCZGqjZ2VGLcpcuNi"),
					entry.WithWebsite("yahoo.com"), entry.WithUsername("swissly0@telemol.online"),
					entry.WithPassword("LxRdhTTSg4Adfkc4"),
				),
				"tEsfGypbPAVCNMAKWtw2mD": entry.New(
					entry.WithId("tEsfGypbPAVCNMAKWtw2mD"),
					entry.WithWebsite("amazon.com"), entry.WithUsername("manhosobpe15@zetgets.com"),
					entry.WithPassword("qHcZsZxPf8acHxxA"),
				),
				"McB6akkM3C5XpzXfMYhasU": entry.New(
					entry.WithId("McB6akkM3C5XpzXfMYhasU"),
					entry.WithWebsite("wikipedia.org"), entry.WithUsername("chuninoleg1971@piftir.com"),
					entry.WithPassword("DT3sftJuRjxWFg68"),
				),
				"YBS32eK8XbeV6ujaY5xERK": entry.New(
					entry.WithId("YBS32eK8XbeV6ujaY5xERK"),
					entry.WithWebsite("twitter.com"), entry.WithUsername("ninablackangel@kumpulanmedia.com"),
					entry.WithPassword("ndUZ6KGduD53up4R"),
				),
				"Fy7HDsbQqkYsbevjuqSG65": entry.New(
					entry.WithId("Fy7HDsbQqkYsbevjuqSG65"),
					entry.WithWebsite("bbc.com"), entry.WithUsername("kgdlove@omdlism.com"),
					entry.WithPassword("jy9EpWExSmmtHa6g"),
				),
				"r5TbidUGZkZeqbP7iCySBn": entry.New(
					entry.WithId("r5TbidUGZkZeqbP7iCySBn"),
					entry.WithWebsite("steampowered.com"), entry.WithUsername("totinoprato@roselarose.com"),
					entry.WithPassword("tbRCJ9uHvxLm9S5q"),
				),
				"pxnChjAmntT5aG35PM3GL4": entry.New(
					entry.WithId("pxnChjAmntT5aG35PM3GL4"),
					entry.WithWebsite("bing.com"), entry.WithUsername("tbiggs@massageshophome.com"),
					entry.WithPassword("H278L5qtwvSVsQzt"),
				),
			},
		),
	)

	s.keystoreRepo.Create(
		keystore.WithName("Work"), keystore.WithEntries(
			map[string]entry.Entry{},
		),
	)

	return s.keystoreRepo.Create(
		keystore.WithName("Other"), keystore.WithEntries(
			map[string]entry.Entry{
				"pxnChjAmntT5aG35PM3G12": entry.New(
					entry.WithId("pxnChjAmntT5aG35PM3G12"),
					entry.WithWebsite("www.microsoft.com"), entry.WithUsername("test123@example.com"),
					entry.WithPassword("H278L5qtwvSVs333"),
				),
			},
		),
	)
}

func (s *service) Keystore(id string) (*keystore.Keystore, error) {
	k, err := s.keystoreRepo.Keystore(id)
	//if errors.Is(err, keystore.ErrAuthenticationRequired) {
	//	return nil, ErrAuthenticationRequired
	//}

	return k, err
}

func (s *service) Keystores() (map[string]*keystore.Keystore, error) {
	ks, err := s.keystoreRepo.Keystores()
	//if err == keystore.ErrAuthenticationRequired {
	//	return nil, ErrAuthenticationRequired
	//} else if err == keystore.ErrInitializationRequired {
	//	return nil, ErrInitializationRequired
	//}

	return ks, err
}

//func (s *service) Unlock(id string, password string) (*keystore.Keystore, error) {
//	k, err := s.keystoreRepo.Unlock(id, password)
//	if errors.Is(err, keystorerepo.ErrAuthenticationFailed) {
//		return nil, ErrAuthenticationFailed
//	}
//
//	return k, err
//}

func (s *service) UpdateKeystore(k *keystore.Keystore) error {
	err := s.keystoreRepo.Update(k)
	//if errors.Is(err, keystore.ErrAuthenticationRequired) {
	//	return ErrAuthenticationRequired
	//}

	return err
}

func (s *service) Authenticate(password string) error {
	return s.keystoreRepo.Authenticate(password)
}

func (s *service) HealthCheck() {
	s.keystoreRepo.HealthCheck()
}

func New(opts ...Option) (keystore.Service, error) {
	s := &service{}

	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	if s.keystoreRepo == nil {
		return nil, ErrInvalidKeystoreRepository
	}

	return s, nil
}
