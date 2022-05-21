package http

import (
	"encoding/hex"

	"myst/internal/server/core/domain/invitation"

	"myst/internal/server/api/http/generated"
	"myst/internal/server/core/domain/keystore"
)

func ToJSONKeystore(k *keystore.Keystore) generated.Keystore {
	return generated.Keystore{
		Id:        k.Id(),
		Name:      k.Name(),
		OwnerId:   k.OwnerId(),
		Payload:   hex.EncodeToString(k.Payload()),
		CreatedAt: k.CreatedAt().Unix(),
		UpdatedAt: k.UpdatedAt().Unix(),
	}
}

func ToJSONInvitation(inv *invitation.Invitation) generated.Invitation {
	gen := generated.Invitation{
		Id:         inv.Id,
		KeystoreId: inv.KeystoreId,
		InviterId:  inv.InviterId,
		InviteeId:  inv.InviteeId,
		CreatedAt:  inv.CreatedAt.Unix(),
		UpdatedAt:  inv.UpdatedAt.Unix(),
	}

	if inv.InviterKey != nil {
		key := inv.InviterKey
		gen.InviterKey = key
	}

	if inv.InviteeKey != nil {
		key := inv.InviteeKey
		gen.InviteeKey = key
	}

	if inv.KeystoreKey != nil {
		key := inv.KeystoreKey
		gen.KeystoreKey = key
	}

	if inv.Finalized() {
		gen.Status = "finalized"
	} else if inv.Accepted() {
		gen.Status = "accepted"
	} else {
		gen.Status = "pending"
	}

	return gen
}
