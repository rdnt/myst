package userservice

import (
	"errors"

	"myst/app/server/domain/keystore"

	"myst/app/server/domain/user"
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

func (s *service) RegisterUser(u *user.User, password string) error {
	panic("implement me")
}

func (s service) CreateKeystore(u *user.User, keystore []byte) error {
	panic("implement me")
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
