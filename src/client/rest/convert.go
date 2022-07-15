package rest

import (
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/user"
	"myst/src/client/rest/generated"
)

func InvitationToRest(inv invitation.Invitation) generated.Invitation {
	return generated.Invitation{
		Id: inv.Id,
		Keystore: generated.Keystore{
			RemoteId: inv.Keystore.RemoteId,
			Name:     inv.Keystore.Name,
		},
		Inviter: generated.User{
			Id:        inv.Inviter.Id,
			Username:  inv.Inviter.Username,
			PublicKey: inv.Inviter.PublicKey,
		},
		Invitee: generated.User{
			Id:        inv.Invitee.Id,
			Username:  inv.Invitee.Username,
			PublicKey: inv.Invitee.PublicKey,
		},
		EncryptedKeystoreKey: inv.EncryptedKeystoreKey,
		Status:               generated.InvitationStatus(inv.Status.String()),
		CreatedAt:            inv.CreatedAt,
		UpdatedAt:            inv.UpdatedAt,
		AcceptedAt:           inv.AcceptedAt,
		DeclinedAt:           inv.DeclinedAt,
		DeletedAt:            inv.DeletedAt,
	}
}

func UserToRest(u user.User) generated.User {
	return generated.User{
		Id:       u.Id,
		Username: u.Username,
	}
}
