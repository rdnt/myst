package userservice

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

func (s *service) Register(opts ...user.Option) (*user.User, error) {
	return s.userRepo.Create(opts...)
}

func (s *service) Authorize(u *user.User, password string) error {
	panic("implement me")
}

func (s *service) User(id string) (*user.User, error) {
	return s.userRepo.User(id)
}

func New(opts ...Option) (user.Service, error) {
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
