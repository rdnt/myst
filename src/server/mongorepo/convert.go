package mongorepo

import (
	"github.com/pkg/errors"

	"myst/src/server/application/domain/invitation"
	"myst/src/server/application/domain/keystore"
	"myst/src/server/application/domain/user"
)

func UserToBSON(u user.User) User {
	return User{
		Id:           u.Id,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		PublicKey:    u.PublicKey,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func UserFromBSON(u User) user.User {
	return user.User{
		Id:           u.Id,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		PublicKey:    u.PublicKey,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func KeystoreToBSON(k keystore.Keystore) Keystore {
	return Keystore{
		Id:        k.Id,
		Name:      k.Name,
		Payload:   k.Payload,
		OwnerId:   k.OwnerId,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func KeystoreFromBSON(k Keystore) keystore.Keystore {
	return keystore.Keystore{
		Id:        k.Id,
		Name:      k.Name,
		Payload:   k.Payload,
		OwnerId:   k.OwnerId,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func InvitationToBSON(inv invitation.Invitation) Invitation {
	return Invitation{
		Id:                   inv.Id,
		KeystoreId:           inv.KeystoreId,
		InviterId:            inv.InviterId,
		InviteeId:            inv.InviteeId,
		EncryptedKeystoreKey: inv.EncryptedKeystoreKey,
		Status:               inv.Status.String(),
		CreatedAt:            inv.CreatedAt,
		UpdatedAt:            inv.UpdatedAt,
		DeletedAt:            inv.DeletedAt,
		DeclinedAt:           inv.DeclinedAt,
		AcceptedAt:           inv.AcceptedAt,
		CancelledAt:          inv.CancelledAt,
		FinalizedAt:          inv.FinalizedAt,
	}
}

func InvitationFromBSON(inv Invitation) (invitation.Invitation, error) {
	stat, err := invitation.StatusFromString(inv.Status)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to parse invitation status")
	}

	return invitation.Invitation{
		Id:                   inv.Id,
		KeystoreId:           inv.KeystoreId,
		InviterId:            inv.InviterId,
		InviteeId:            inv.InviteeId,
		EncryptedKeystoreKey: inv.EncryptedKeystoreKey,
		Status:               stat,
		CreatedAt:            inv.CreatedAt,
		UpdatedAt:            inv.UpdatedAt,
		DeletedAt:            inv.DeletedAt,
		DeclinedAt:           inv.DeclinedAt,
		AcceptedAt:           inv.AcceptedAt,
		CancelledAt:          inv.CancelledAt,
		FinalizedAt:          inv.FinalizedAt,
	}, nil
}
