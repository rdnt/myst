package http

import (
	"myst/internal/server/api/http/generated"
	"myst/internal/server/core/domain/invitation"
	"myst/internal/server/core/domain/keystore"
	"myst/internal/server/core/domain/user"
)

func ToJSONKeystore(k *keystore.Keystore) generated.Keystore {
	return generated.Keystore{
		Id:        k.Id,
		Name:      k.Name,
		OwnerId:   k.OwnerId,
		Payload:   k.Payload,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func ToJSONInvitation(inv invitation.Invitation) generated.Invitation {
	return generated.Invitation{
		Id:           inv.Id,
		KeystoreId:   inv.KeystoreId,
		KeystoreName: inv.KeystoreName,
		InviterId:    inv.InviterId,
		InviteeId:    inv.InviteeId,
		Status:       generated.InvitationStatus(inv.Status.String()),
		InviterKey:   inv.InviterKey,
		InviteeKey:   inv.InviteeKey,
		KeystoreKey:  inv.KeystoreKey,
		CreatedAt:    inv.CreatedAt,
		UpdatedAt:    inv.UpdatedAt,
		AcceptedAt:   inv.AcceptedAt,
		DeclinedAt:   inv.DeclinedAt,
		DeletedAt:    inv.DeletedAt,
	}
}

func ToJSONUser(u user.User) generated.User {
	return generated.User{
		Id:        u.Id,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
