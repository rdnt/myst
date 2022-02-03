package keystoreservice

import (
	"errors"
	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/application/domain/keystore/entry"
	"myst/pkg/logger"
)

var (
	ErrInvalidKeystoreRepository = errors.New("invalid keystore repository")
	ErrAuthenticationRequired    = errors.New("authentication required")
	ErrAuthenticationFailed      = errors.New("authentication failed")
	ErrInitializationRequired    = errors.New("initialization required")
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

func (s *service) Create(name string) (*keystore.Keystore, error) {
	return s.keystoreRepo.Create(keystore.WithName(name))
}

func (s *service) Initialize(name, password string) (*keystore.Keystore, error) {
	err := s.keystoreRepo.Initialize(password)
	if err != nil {
		return nil, err
	}

	// TODO: remove dummy keystores and properly return error
	s.keystoreRepo.Create(
		keystore.WithName(name), keystore.WithEntries(
			[]entry.Entry{
				entry.New(
					entry.WithWebsite("github.com"), entry.WithUsername("rdntdev@gmail.com"),
					entry.WithPassword("nzK&d#u+MjFU8p&4UhL)s3+h"),
					entry.WithNotes("Lorem ipsum"),
				),
				entry.New(
					entry.WithWebsite("youtube.com"), entry.WithUsername("oldsnut@gmailni.com"),
					entry.WithPassword("tsksWgABXhvh9LfF"),
					entry.WithNotes("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."),
				),
				entry.New(
					entry.WithWebsite("facebook.com"), entry.WithUsername("pete24uk@test130.com"),
					entry.WithPassword("uXekxDRk6bmvvpda"),
					entry.WithNotes("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."),
				),
				entry.New(
					entry.WithWebsite("baidu.com"), entry.WithUsername("swissly0@telemol.online"),
					entry.WithPassword("6Yu4k2YkNPxpZkHn"),
				),
				entry.New(
					entry.WithWebsite("yahoo.com"), entry.WithUsername("swissly0@telemol.online"),
					entry.WithPassword("LxRdhTTSg4Adfkc4"),
				),
				entry.New(
					entry.WithWebsite("amazon.com"), entry.WithUsername("manhosobpe15@zetgets.com"),
					entry.WithPassword("qHcZsZxPf8acHxxA"),
				),
				entry.New(
					entry.WithWebsite("wikipedia.org"), entry.WithUsername("chuninoleg1971@piftir.com"),
					entry.WithPassword("DT3sftJuRjxWFg68"),
				),
				entry.New(
					entry.WithWebsite("twitter.com"), entry.WithUsername("ninablackangel@kumpulanmedia.com"),
					entry.WithPassword("ndUZ6KGduD53up4R"),
				),
				entry.New(
					entry.WithWebsite("bbc.com"), entry.WithUsername("kgdlove@omdlism.com"),
					entry.WithPassword("jy9EpWExSmmtHa6g"),
				),
				entry.New(
					entry.WithWebsite("steampowered.com"), entry.WithUsername("totinoprato@roselarose.com"),
					entry.WithPassword("tbRCJ9uHvxLm9S5q"),
				),
				entry.New(
					entry.WithWebsite("bing.com"), entry.WithUsername("tbiggs@massageshophome.com"),
					entry.WithPassword("H278L5qtwvSVsQzt"),
				),
			},
		),
	)

	s.keystoreRepo.Create(
		keystore.WithName("Work"), keystore.WithEntries(
			[]entry.Entry{},
		),
	)

	return s.keystoreRepo.Create(
		keystore.WithName("Other"), keystore.WithEntries(
			[]entry.Entry{},
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

func (s *service) Update(k *keystore.Keystore) error {
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
