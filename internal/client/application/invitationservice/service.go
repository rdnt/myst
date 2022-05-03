package invitationservice

import (
	"myst/internal/client/application/domain/invitation"
	"myst/pkg/logger"

	"github.com/pkg/errors"
)

var (
	ErrInvalidInvitationsRepository = errors.New("invalid invitations repository")
)

type service struct {
	invitations invitation.Repository
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

	if s.invitations == nil {
		return nil, ErrInvalidInvitationsRepository
	}

	return s, nil
}

func (s *service) CreateKeystoreInvitation(keystoreId, inviterId, inviteeId string, inviterKey []byte) (*invitation.Invitation, error) {
	//TODO implement me
	panic("implement me")
}
