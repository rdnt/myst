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

	return s.keystoreRepo.Create(
		keystore.WithName(name), keystore.WithEntries(
			[]entry.Entry{
				entry.New(entry.WithLabel("google.com"), entry.WithUsername("rdnt"), entry.WithPassword("1234")),
				entry.New(entry.WithLabel("google.com"), entry.WithUsername("rdnt"), entry.WithPassword("1234")),
				entry.New(entry.WithLabel("google.com"), entry.WithUsername("rdnt"), entry.WithPassword("1234")),
				entry.New(entry.WithLabel("google.com"), entry.WithUsername("rdnt"), entry.WithPassword("1234")),
				entry.New(entry.WithLabel("google.com"), entry.WithUsername("rdnt"), entry.WithPassword("1234")),
				entry.New(entry.WithLabel("google.com"), entry.WithUsername("rdnt"), entry.WithPassword("1234")),
				entry.New(entry.WithLabel("google.com"), entry.WithUsername("rdnt"), entry.WithPassword("1234")),
				entry.New(entry.WithLabel("google.com"), entry.WithUsername("rdnt"), entry.WithPassword("1234")),
				entry.New(entry.WithLabel("google.com"), entry.WithUsername("rdnt"), entry.WithPassword("1234")),
				entry.New(entry.WithLabel("google.com"), entry.WithUsername("rdnt"), entry.WithPassword("1234")),
				entry.New(entry.WithLabel("google.com"), entry.WithUsername("rdnt"), entry.WithPassword("1234")),
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
