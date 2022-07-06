package application

import (
	"github.com/pkg/errors"

	"myst/internal/server/core/domain/invitation"
)

func (app *application) CreateInvitation(keystoreId, inviterId, inviteeId string) (invitation.Invitation, error) {
	if inviterId == inviteeId {
		return invitation.Invitation{}, errors.New("inviter cannot be the same as invitee")
	}

	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	inviter, err := app.users.User(inviterId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	invitee, err := app.users.User(inviteeId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	return app.invitations.CreateInvitation(
		invitation.WithKeystoreId(k.Id),
		invitation.WithInviterId(inviter.Id),
		invitation.WithInviteeId(invitee.Id),
	)
}

func (app *application) AcceptInvitation(userId string, invitationId string) (invitation.Invitation, error) {
	inv, err := app.UserInvitation(userId, invitationId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	err = inv.Accept()
	if err != nil {
		return invitation.Invitation{}, err
	}

	err = app.invitations.UpdateInvitation(&inv)
	if err != nil {
		return invitation.Invitation{}, err
	}

	return inv, nil
}

func (app *application) DeclineOrCancelInvitation(userId, invitationId string) (invitation.Invitation, error) {
	u, err := app.users.User(userId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	inv, err := app.invitations.Invitation(invitationId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	if u.Id != inv.InviterId && u.Id != inv.InviteeId {
		return invitation.Invitation{}, errors.New("unauthorized")
	}

	if u.Id == inv.InviterId {
		err = inv.Delete()
		if err != nil {
			return invitation.Invitation{}, err
		}
	} else if u.Id == inv.InviteeId {
		err = inv.Decline()
		if err != nil {
			return invitation.Invitation{}, err
		}
	}

	err = app.invitations.UpdateInvitation(&inv)
	if err != nil {
		return invitation.Invitation{}, err
	}

	return inv, nil
}

func (app *application) FinalizeInvitation(invitationId string, encryptedKeystoreKey []byte) (invitation.Invitation, error) {
	inv, err := app.invitations.Invitation(invitationId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	err = inv.Finalize(encryptedKeystoreKey)
	if err != nil {
		return invitation.Invitation{}, err
	}

	err = app.invitations.UpdateInvitation(&inv)
	if err != nil {
		return invitation.Invitation{}, err
	}

	return inv, nil
}
