package invitationservice

import (
	"errors"

	"myst/internal/server/core/domain/invitation"
	"myst/internal/server/core/domain/keystore"
	"myst/internal/server/core/domain/user"

	"myst/pkg/logger"
)

var (
	ErrInvalidKeystoreRepository = errors.New("invalid keystore repository")
	ErrInvalidUserRepository     = errors.New("invalid user repository")
)

var log = logger.New("invitations", logger.Green)

type service struct {
	userRepo       user.Repository
	keystoreRepo   keystore.Repository
	invitationRepo invitation.Repository
}

func (s *service) Create(opts ...invitation.Option) (*invitation.Invitation, error) {
	// TODO

	log.Debug("CreateKeystoreInvitation", opts)

	inv, err := s.invitationRepo.Create(opts...)
	if err != nil {
		return nil, err
	}

	log.Debug(inv)

	return inv, nil
}

func New(opts ...Option) (invitation.Service, error) {
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

	if s.invitationRepo == nil {
		return nil, ErrInvalidUserRepository
	}

	return s, nil
}
