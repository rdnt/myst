package invitationservice

import (
	"myst/internal/server/core/domain/invitation"
	"myst/internal/server/core/domain/keystore"
	"myst/internal/server/core/domain/user"
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

func WithInvitationRepository(repo invitation.Repository) Option {
	return func(s *service) error {
		s.invitationRepo = repo
		return nil
	}
}
