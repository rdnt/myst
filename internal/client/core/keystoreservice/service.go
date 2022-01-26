package keystoreservice

import (
	"errors"

	application "myst/internal/client"
	"myst/internal/client/core/domain/keystore"
	keystorerepo "myst/internal/client/core/keystorerepo/fs"

	"myst/pkg/logger"
)

var (
	ErrInvalidKeystoreRepository = errors.New("invalid keystore repository")
	ErrAuthenticationRequired    = keystore.ErrAuthenticationRequired
	ErrAuthenticationFailed      = keystorerepo.ErrAuthenticationFailed
)

type service struct {
	keystoreRepo application.KeystoreRepository
}

func (s *service) KeystoreIds() ([]string, error) {
	return s.keystoreRepo.KeystoreIds()
}

func (s *service) Create(opts ...keystore.Option) (*keystore.Keystore, error) {
	return s.keystoreRepo.Create(opts...)
}

func (s *service) Keystore(id string) (*keystore.Keystore, error) {
	k, err := s.keystoreRepo.Keystore(id)
	if errors.Is(err, keystore.ErrAuthenticationRequired) {
		return nil, ErrAuthenticationRequired
	}

	return k, err
}

func (s *service) Keystores() ([]*keystore.Keystore, error) {
	ks, err := s.keystoreRepo.Keystores()
	if errors.Is(err, keystore.ErrAuthenticationRequired) {
		return nil, ErrAuthenticationRequired
	}

	return ks, err
}

func (s *service) Unlock(id string, password string) (*keystore.Keystore, error) {
	k, err := s.keystoreRepo.Unlock(id, password)
	if errors.Is(err, keystorerepo.ErrAuthenticationFailed) {
		return nil, ErrAuthenticationFailed
	}

	return k, err
}

func (s *service) Update(k *keystore.Keystore) error {
	err := s.keystoreRepo.Update(k)
	if errors.Is(err, keystore.ErrAuthenticationRequired) {
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
