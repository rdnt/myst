package keystoreservice

import (
	"errors"

	"myst/app/client/core/domain/keystore"
	"myst/pkg/logger"
)

var (
	ErrInvalidKeystoreRepository = errors.New("invalid keystore repository")
)

type service struct {
	keystoreRepo keystore.Repository
}

func (s *service) Create(opts ...keystore.Option) (*keystore.Keystore, error) {
	k, err := s.keystoreRepo.Create(opts...)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (s *service) Keystore(id string) (*keystore.Keystore, error) {
	return s.keystoreRepo.Keystore(id)
}

func (s *service) Update(k *keystore.Keystore) error {
	return s.keystoreRepo.Update(k)
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
