package http

import (
	"myst/internal/client/api/http/generated"
	"myst/internal/client/application/domain/invitation"
)

func InvitationToRest(inv invitation.Invitation) generated.Invitation {
	return generated.Invitation{
		Id:           inv.Id,
		KeystoreId:   inv.KeystoreId,
		KeystoreName: inv.KeystoreName,
		InviterId:    inv.InviterId,
		InviteeId:    inv.InviteeId,
		InviterKey:   inv.InviterKey,
		InviteeKey:   inv.InviteeKey,
		KeystoreKey:  inv.KeystoreKey,
		Status:       generated.InvitationStatus(inv.Status.String()),
		CreatedAt:    inv.CreatedAt.Unix(),
		UpdatedAt:    inv.UpdatedAt.Unix(),
	}
}
