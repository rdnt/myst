package rest

import (
	"myst/pkg/hashicon"
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
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
		Inviter:              UserToRest(inv.Inviter),
		Invitee:              UserToRest(inv.Invitee),
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
	h, err := hashicon.New(u.PublicKey)
	if err != nil {
		panic(err)
	}

	return generated.User{
		Id:        u.Id,
		Username:  u.Username,
		PublicKey: u.PublicKey,
		Icon:      h.ToSVG(),
	}
}

func KeystoreToRest(k keystore.Keystore) generated.Keystore {
	entries := []generated.Entry{}

	for _, e := range k.Entries {
		entries = append(
			entries, generated.Entry{
				Id:       e.Id,
				Website:  e.Website,
				Username: e.Username,
				Password: e.Password,
				Notes:    e.Notes,
			},
		)
	}

	return generated.Keystore{
		Id:       k.Id,
		RemoteId: k.RemoteId,
		ReadOnly: k.ReadOnly,
		Name:     k.Name,
		Entries:  entries,
	}
}
