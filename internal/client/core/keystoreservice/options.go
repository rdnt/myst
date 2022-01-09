package keystoreservice

import (
	"myst/internal/client"
)

type Option func(s *service) error

func WithKeystoreRepository(repo client.KeystoreRepository) Option {
	return func(s *service) error {
		s.keystoreRepo = repo
		return nil
	}
}
