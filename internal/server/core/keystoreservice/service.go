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

	return s.keystoreRepo.Create(
		keystore.WithName(name),
		keystore.WithOwnerId(u.Id()),
		keystore.WithPayload(payload),
	)
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