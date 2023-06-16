package rest

import (
	"github.com/pkg/errors"

	"myst/src/server/application/domain/invitation"
	"myst/src/server/application/domain/keystore"
	"myst/src/server/application/domain/user"
	"myst/src/server/rest/generated"
)

func keystoreToJSON(k keystore.Keystore) generated.Keystore {
	return generated.Keystore{
		Id:        k.Id,
		Name:      k.Name,
		OwnerId:   k.OwnerId,
		Payload:   k.Payload,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func (s *Server) invitationToJSON(inv invitation.Invitation) (generated.Invitation, error) {
	inviter, err := s.app.User(inv.InviterId)
	if err != nil {
		return generated.Invitation{}, errors.WithMessage(err, "failed to get inviter")
	}

	invitee, err := s.app.User(inv.InviteeId)
	if err != nil {
		return generated.Invitation{}, errors.WithMessage(err, "failed to get invitee")
	}

	k, err := s.app.Keystore(inv.KeystoreId)
	if err != nil {
		return generated.Invitation{}, errors.WithMessage(err, "failed to get keystore")
	}

	return generated.Invitation{
		Id: inv.Id,
		Keystore: generated.KeystoreName{
			Id:   k.Id,
			Name: k.Name,
		},
		Inviter:              userToJSON(inviter),
		Invitee:              userToJSON(invitee),
		Status:               generated.InvitationStatus(inv.Status.String()),
		EncryptedKeystoreKey: inv.EncryptedKeystoreKey,
		CreatedAt:            inv.CreatedAt,
		UpdatedAt:            inv.UpdatedAt,
		DeletedAt:            inv.DeletedAt,
		DeclinedAt:           inv.DeclinedAt,
		AcceptedAt:           inv.AcceptedAt,
		CancelledAt:          inv.CancelledAt,
		FinalizedAt:          inv.FinalizedAt,
	}, nil
}

func userToJSON(u user.User) generated.User {
	return generated.User{
		Id:        u.Id,
		PublicKey: u.PublicKey,
		Username:  u.Username,
	}
}
