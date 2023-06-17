package rest

import (
	"myst/pkg/hashicon"
	"myst/pkg/optional"
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/user"
	"myst/src/client/rest/generated"
)

func (s *Server) invitationToJSON(inv invitation.Invitation) (generated.Invitation, error) {
	inviter, err := s.userToJSON(inv.Inviter)
	if err != nil {
		return generated.Invitation{}, err
	}

	invitee, err := s.userToJSON(inv.Invitee)
	if err != nil {
		return generated.Invitation{}, err
	}

	return generated.Invitation{
		Id: inv.Id,
		Keystore: generated.InvitationKeystore{
			RemoteId: inv.RemoteKeystoreId,
			Name:     inv.KeystoreName,
		},
		Inviter:              inviter,
		Invitee:              invitee,
		EncryptedKeystoreKey: inv.EncryptedKeystoreKey,
		Status:               generated.InvitationStatus(inv.Status.String()),
		CreatedAt:            inv.CreatedAt,
		UpdatedAt:            inv.UpdatedAt,
		AcceptedAt:           inv.AcceptedAt,
		DeclinedAt:           inv.DeclinedAt,
		DeletedAt:            inv.DeletedAt,
	}, nil
}

func (s *Server) userToJSON(u user.User) (generated.User, error) {
	var icon *string

	sharedSecret, err := s.app.SharedSecret(u.Id)
	if err != nil {
		return generated.User{}, err
	}

	ic, err := hashicon.New(sharedSecret)
	if err != nil {
		return generated.User{}, err
	}

	icon = optional.Ref(ic.ToSVG())

	return generated.User{
		Id:        u.Id,
		Username:  u.Username,
		PublicKey: u.PublicKey,
		Icon:      icon,
	}, nil
}

func keystoreToJSON(k keystore.Keystore) generated.Keystore {
	entries := make([]generated.Entry, 0)

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
