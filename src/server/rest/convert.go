package rest

import (
	"myst/src/server/application/domain/invitation"
	"myst/src/server/application/domain/keystore"
	"myst/src/server/application/domain/user"
	"myst/src/server/rest/generated"
)

func ToJSONKeystore(k keystore.Keystore) generated.Keystore {
	return generated.Keystore{
		Id:        k.Id,
		Name:      k.Name,
		OwnerId:   k.OwnerId,
		Payload:   k.Payload,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func (s *Server) ToJSONInvitation(i invitation.Invitation) (generated.Invitation, error) {
	inviter, err := s.app.User(i.InviterId)
	if err != nil {
		return generated.Invitation{}, err
	}

	invitee, err := s.app.User(i.InviteeId)
	if err != nil {
		return generated.Invitation{}, err
	}

	k, err := s.app.Keystore(i.KeystoreId)
	if err != nil {
		return generated.Invitation{}, err
	}

	return generated.Invitation{
		Id: i.Id,
		Keystore: generated.KeystoreName{
			Id:   k.Id,
			Name: k.Name,
		},
		Inviter: generated.User{
			Id:        inviter.Id,
			Username:  inviter.Username,
			PublicKey: inviter.PublicKey,
		},
		Invitee: generated.User{
			Id:        invitee.Id,
			Username:  invitee.Username,
			PublicKey: invitee.PublicKey,
		},
		Status:               generated.InvitationStatus(i.Status.String()),
		EncryptedKeystoreKey: i.EncryptedKeystoreKey,
		CreatedAt:            i.CreatedAt,
		UpdatedAt:            i.UpdatedAt,
		AcceptedAt:           i.AcceptedAt,
		DeclinedAt:           i.DeclinedAt,
		DeletedAt:            i.DeletedAt,
	}, nil
}

func UserToJSON(u user.User) generated.User {
	return generated.User{
		Id:       u.Id,
		Username: u.Username,
	}
}
