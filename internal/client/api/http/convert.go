package http

import (
	"myst/internal/client/api/http/generated"
	"myst/internal/client/application/domain/invitation"
	"myst/internal/client/application/domain/user"
)

func InvitationToRest(inv invitation.Invitation) generated.Invitation {
	return generated.Invitation{
		Id:                   inv.Id,
		KeystoreId:           inv.KeystoreId,
		KeystoreName:         inv.KeystoreName,
		InviterId:            inv.InviterId,
		InviteeId:            inv.InviteeId,
		InviterPublicKey:     inv.InviterPublicKey,
		InviteePublicKey:     inv.InviteePublicKey,
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
