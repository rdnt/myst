package keystoreservice

import (
	"errors"

	"myst/internal/server/core/domain/keystore"
	"myst/internal/server/core/domain/user"

	"myst/pkg/logger"
)

var (
	ErrInvalidKeystoreRepository = errors.New("invalid keystore repository")
	ErrInvalidUserRepository     = errors.New("invalid user repository")
)

type service struct {
	userRepo     user.Repository
	keystoreRepo keystore.Repository
}

func (s *service) Create(name, ownerId string, payload []byte) (*keystore.Keystore, error) {
	u, err := s.userRepo.User(ownerId)
	if err != nil {
		return nil, err
	}

	k, err := s.keystoreRepo.Create(
		keystore.WithName(name),
		keystore.WithOwnerId(u.Id()),
		keystore.WithPayload(payload),
	)

	u.OwnKeystore(k.Id())

	err = s.userRepo.Update(u)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (s *service) Keystore(id string) (*keystore.Keystore, error) {
	return s.keystoreRepo.Keystore(id)
}

func (s *service) Keystores() ([]*keystore.Keystore, error) {
	return s.keystoreRepo.Keystores()
}

func (s *service) UserKeystore(userId, keystoreId string) (*keystore.Keystore, error) {
	u, err := s.userRepo.User(userId)
	if err != nil {
		return nil, err
	}

	var exists bool
	for _, kid := range u.KeystoreIds() {
		if kid == keystoreId {
			exists = true
			break
		}
	}

	if !exists {
		return nil, keystore.ErrNotFound
	}

	return s.keystoreRepo.Keystore(keystoreId)
}

func (s *service) UserKeystores(userId string) ([]*keystore.Keystore, error) {
	u, err := s.userRepo.User(userId)
	if err != nil {
		return nil, err
	}

	ks := []*keystore.Keystore{}
	for _, kid := range u.KeystoreIds() {
		k, err := s.keystoreRepo.Keystore(kid)
		if err != nil {
			return nil, err
		}

		ks = append(ks, k)
	}

	return ks, nil
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

	if s.userRepo == nil {
		return nil, ErrInvalidUserRepository
	}

	return s, nil
}
