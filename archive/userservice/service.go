package userservice

import (
	"errors"

	"myst/internal/server/applicationrefactor/domain/device"

	"myst/internal/server/application/domain/keystore"
	"myst/internal/server/application/domain/user"
	"myst/pkg/logger"
)

var (
	ErrInvalidKeystoreRepository = errors.New("invalid keystore repository")
	ErrInvalidUserRepository     = errors.New("invalid user repository")
	ErrInvalidDeviceRepository   = errors.New("invalid device repository")
)

type service struct {
	userRepo     user.Repository
	keystoreRepo keystore.Repository
	deviceRepo   device.Repository
}

func (s *service) UserDevices(userId string) (map[string]device.Device, error) {
	return s.deviceRepo.UserDevices(userId)
}

func (s *service) CreateUser(opts ...user.Option) (user.User, error) {
	return s.userRepo.CreateUser(opts...)
}

func (s *service) AuthorizeUser(u user.User, password string) error {
	panic("implement me")
}

func (s *service) User(id string) (user.User, error) {
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

	if s.deviceRepo == nil {
		return nil, ErrInvalidDeviceRepository
	}

	return s, nil
}
