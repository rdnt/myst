package invitationservice

import (
	"myst/internal/server/core/domain/invitation"
	"myst/internal/server/core/domain/keystore"
	"myst/internal/server/core/domain/user"
	"myst/pkg/logger"

	"github.com/pkg/errors"
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

func (s *service) CreateInvitation(keystoreId, inviterId, inviteeId string, inviterKey []byte) (invitation.Invitation, error) {
	if inviterId == inviteeId {
		return invitation.Invitation{}, errors.New("inviter cannot be the same as invitee")
	}

	k, err := s.keystoreRepo.Keystore(keystoreId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	//inviter, err := s.userRepo.User(inviterId)
	//if err != nil {
	//	return nil, err
	//}
	//
	//invitee, err := s.userRepo.User(inviteeId)
	//if err != nil {
	//	return nil, err
	//}

	return s.invitationRepo.CreateInvitation(
		invitation.WithKeystoreId(k.Id),
		invitation.WithInviterId(inviterId),
		invitation.WithInviteeId(inviteeId),
	)
}

func (s *service) Invitation(id string) (invitation.Invitation, error) {
	return s.invitationRepo.Invitation(id)
}

func (s *service) AcceptInvitation(invitationId string, inviteeKey []byte) (invitation.Invitation, error) {
	inv, err := s.invitationRepo.Invitation(invitationId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	err = inv.Accept(inviteeKey)
	if err != nil {
		return invitation.Invitation{}, err
	}

	err = s.invitationRepo.UpdateInvitation(&inv)
	if err != nil {
		return invitation.Invitation{}, err
	}

	return inv, nil
}

func (s *service) DeclineOrCancelInvitation(userId, invitationId string) (invitation.Invitation, error) {
	inv, err := s.invitationRepo.Invitation(invitationId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	if userId != inv.InviterId && userId != inv.InviteeId {
		return invitation.Invitation{}, errors.New("unauthorized")
	}

	if userId == inv.InviterId {
		err = inv.Delete()
		if err != nil {
			return invitation.Invitation{}, err
		}
	} else if userId == inv.InviteeId {
		err = inv.Decline()
		if err != nil {
			return invitation.Invitation{}, err
		}
	}

	err = s.invitationRepo.UpdateInvitation(&inv)
	if err != nil {
		return invitation.Invitation{}, err
	}

	return inv, nil
}

func (s *service) FinalizeInvitation(invitationId string, keystoreKey []byte) (invitation.Invitation, error) {
	inv, err := s.invitationRepo.Invitation(invitationId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	err = inv.Finalize(keystoreKey)
	if err != nil {
		return invitation.Invitation{}, err
	}

	err = s.invitationRepo.UpdateInvitation(&inv)
	if err != nil {
		return invitation.Invitation{}, err
	}

	return inv, nil
}

type UserInvitationsOptions struct {
	Status *string
}

// UserInvitations returns all the invitations this user has access to. These include:
// - invitations where the user is the inviter
// - invitations where the user is the invitee
func (s *service) UserInvitations(deviceId string, opts *invitation.UserInvitationsOptions) ([]invitation.Invitation, error) {
	//u, err := s.userRepo.User(userId)
	//if err != nil {
	//	return nil, err
	//}

	invs, err := s.invitationRepo.UserInvitations(deviceId)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get user invitations")
	}

	invitations := []invitation.Invitation{}
	for _, inv := range invs {
		if inv.Deleted() && inv.InviteeId == deviceId {
			continue
		}
		if opts != nil && opts.Status != nil && *opts.Status != inv.Status {
			continue
		}

		invitations = append(invitations, inv)
	}

	return invitations, nil
}

// UserInvitation returns an invitation that a user has access to.
func (s *service) UserInvitation(userId, invitationId string) (invitation.Invitation, error) {
	return s.invitationRepo.UserInvitation(userId, invitationId)
}

//func (s *service) UserKeystores(userId string) ([]*keystore.Keystore, error) {
//	invs, err := s.UserInvitations(userId)
//	if err != nil {
//		return nil, errors.WithMessage(err, "failed to get user invitations")
//	}
//
//	keystores := []*keystore.Keystore{}
//
//	for _, inv := range invs {
//		// TODO: split logic to separate function "AcceptedUserInviations"
//		if inv.Finalized() && inv.InviteeId() == userId {
//			k, err := s.keystoreRepo.Keystore(inv.KeystoreId())
//			if err != nil {
//				return nil, errors.WithMessage(err, "failed to get keystore")
//			}
//
//			keystores = append(keystores, k)
//		}
//	}
//
//	u, err := s.userRepo.CurrentUser(userId)
//	if err != nil {
//		return nil, err
//	}
//
//	ks, err := s.keystoreRepo.UserKeystores(u.Id())
//	if err != nil {
//		return nil, errors.WithMessage(err, "failed to get user keystores")
//	}
//
//	return append(keystores, ks...), nil
//}

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
