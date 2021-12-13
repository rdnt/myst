package keystoreservice

import (
	"myst/app/client/core/domain/keystore"
)

type Option func(s *service) error

func WithKeystoreRepository(repo keystore.Repository) Option {
	return func(s *service) error {
		s.keystoreRepo = repo
		return nil
	}
}
