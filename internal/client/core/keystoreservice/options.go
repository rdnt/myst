package keystoreservice

import (
	application "myst/internal/client"
)

type Option func(s *service) error

func WithKeystoreRepository(repo application.KeystoreRepository) Option {
	return func(s *service) error {
		s.keystoreRepo = repo
		return nil
	}
}
