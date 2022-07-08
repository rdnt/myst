package rest

import (
	"myst/internal/server/application/domain/invitation"
	"myst/internal/server/application/domain/keystore"
	"myst/internal/server/application/domain/user"
	"myst/internal/server/rest/generated"
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

func (api *Server) ToJSONInvitation(i invitation.Invitation) (generated.Invitation, error) {
	inviter, err := api.app.User(i.InviterId)
	if err != nil {
		return generated.Invitation{}, err
	}

	invitee, err := api.app.User(i.InviteeId)
	if err != nil {
		return generated.Invitation{}, err
	}

	k, err := api.app.Keystore(i.KeystoreId)
	if err != nil {
		return generated.Invitation{}, err
	}

	return generated.Invitation{
		Id:                   i.Id,
		KeystoreId:           k.Id,
		KeystoreName:         k.Name,
		InviterId:            inviter.Id,
		InviterPublicKey:     inviter.PublicKey,
		InviteeId:            invitee.Id,
		InviteePublicKey:     invitee.PublicKey,
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
