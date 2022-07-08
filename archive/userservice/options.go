package userservice

import (
	"myst/internal/server/applicationrefactor/domain/device"

	"myst/internal/server/application/domain/keystore"
	"myst/internal/server/application/domain/user"
)

type Option func(s *service) error

func WithUserRepository(repo user.Repository) Option {
	return func(s *service) error {
		s.userRepo = repo
		return nil
	}
}

func WithKeystoreRepository(repo keystore.Repository) Option {
	return func(s *service) error {
		s.keystoreRepo = repo
		return nil
	}
}

func WithDeviceRepository(repo device.Repository) Option {
	return func(s *service) error {
		s.deviceRepo = repo
		return nil
	}
}
