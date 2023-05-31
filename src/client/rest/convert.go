package rest

import (
	"myst/pkg/hashicon"
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/user"
	"myst/src/client/rest/generated"
)

func (s *Server) InvitationToRest(inv invitation.Invitation) (generated.Invitation, error) {
	inviter, err := s.userToRest(inv.Inviter)
	if err != nil {
		return generated.Invitation{}, err
	}

	invitee, err := s.userToRest(inv.Invitee)
	if err != nil {
		return generated.Invitation{}, err
	}

	return generated.Invitation{
		Id: inv.Id,
		Keystore: generated.Keystore{
			RemoteId: inv.Keystore.RemoteId,
			Name:     inv.Keystore.Name,
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

func (s *Server) userToRest(u user.User) (generated.User, error) {
	var icon *string
	if u.SharedSecret != nil {
		ic, err := hashicon.New(u.SharedSecret)
		if err != nil {
			return generated.User{}, err
		}

		str := ic.ToSVG()
		icon = &str
	}

	return generated.User{
		Id:        u.Id,
		Username:  u.Username,
		PublicKey: u.PublicKey,
		Icon:      icon,
	}, nil
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
