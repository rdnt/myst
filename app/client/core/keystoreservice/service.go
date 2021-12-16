package keystoreservice

import (
	"errors"

	jsonkeystore "myst/app/client/core/keystorerepo/keystore"

	"myst/app/client"

	"myst/app/client/core/domain/keystore"
	"myst/pkg/logger"
)

var (
	ErrInvalidKeystoreRepository = errors.New("invalid keystore repository")
	ErrAuthenticationRequired    = errors.New("authentication required")
)

type service struct {
	keystoreRepo client.KeystoreRepository
}

func (s *service) Create(opts ...keystore.Option) (*keystore.Keystore, error) {
	return s.keystoreRepo.Create(opts...)
}

func (s *service) Keystore(id string) (*keystore.Keystore, error) {
	k, err := s.keystoreRepo.Keystore(id)
	if errors.Is(err, jsonkeystore.ErrAuthenticationRequired) {
		return nil, ErrAuthenticationRequired
	}

	return k, err
}

func (s *service) Unlock(id string, passphrase string) (*keystore.Keystore, error) {
	return s.keystoreRepo.Unlock(id, passphrase)
}

func (s *service) Update(k *keystore.Keystore) error {
	err := s.keystoreRepo.Update(k)
	if errors.Is(err, jsonkeystore.ErrAuthenticationRequired) {
		return ErrAuthenticationRequired
	}

	return err
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
