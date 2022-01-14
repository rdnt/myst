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

func (s *service) Create(keystoreId, inviterId, inviteeId string, inviterKey []byte) (*invitation.Invitation, error) {
	store, err := s.keystoreRepo.Keystore(keystoreId)
	if err != nil {
		return nil, err
	}

	inviter, err := s.userRepo.User(inviterId)
	if err != nil {
		return nil, err
	}

	invitee, err := s.userRepo.User(inviteeId)
	if err != nil {
		return nil, err
	}

	return s.invitationRepo.Create(
		invitation.WithKeystore(store),
		invitation.WithInviter(inviter),
		invitation.WithInvitee(invitee),
		invitation.WithInviterKey(inviterKey),
	)
}

func (s *service) Accept(invitationId string, inviteeKey []byte) (*invitation.Invitation, error) {
	inv, err := s.invitationRepo.Invitation(invitationId)
	if err != nil {
		return nil, err
	}

	err = inv.Accept(inviteeKey)
	if err != nil {
		return nil, err
	}

	err = s.invitationRepo.Update(inv)
	if err != nil {
		return nil, err
	}

	return inv, nil
}

func (s *service) Finalize(invitationId string, keystoreKey []byte) (*invitation.Invitation, error) {
	inv, err := s.invitationRepo.Invitation(invitationId)
	if err != nil {
		return nil, err
	}

	err = inv.Finalize(keystoreKey)
	if err != nil {
		return nil, err
	}

	err = s.invitationRepo.Update(inv)
	if err != nil {
		return nil, err
	}

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
