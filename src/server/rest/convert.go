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

func (s *Server) ToJSONInvitation(inv invitation.Invitation) (generated.Invitation, error) {
	inviter, err := s.app.User(inv.InviterId)
	if err != nil {
		return generated.Invitation{}, err
	}

	invitee, err := s.app.User(inv.InviteeId)
	if err != nil {
		return generated.Invitation{}, err
	}

	k, err := s.app.Keystore(inv.KeystoreId)
	if err != nil {
		return generated.Invitation{}, err
	}

	return generated.Invitation{
		Id: inv.Id,
		Keystore: generated.KeystoreName{
			Id:   k.Id,
			Name: k.Name,
		},
		Inviter:              UserToJSON(inviter),
		Invitee:              UserToJSON(invitee),
		Status:               generated.InvitationStatus(inv.Status.String()),
		EncryptedKeystoreKey: inv.EncryptedKeystoreKey,
		CreatedAt:            inv.CreatedAt,
		UpdatedAt:            inv.UpdatedAt,
		AcceptedAt:           inv.AcceptedAt,
		DeclinedAt:           inv.DeclinedAt,
		DeletedAt:            inv.DeletedAt,
	}, nil
}

func UserToJSON(u user.User) generated.User {
	return generated.User{
		Id:        u.Id,
		PublicKey: u.PublicKey,
		Username:  u.Username,
	}
}
