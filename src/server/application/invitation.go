package application

import (
	"github.com/pkg/errors"

	"myst/src/server/application/domain/invitation"
)

func (app *application) CreateInvitation(keystoreId, inviterId, inviteeUsername string) (invitation.Invitation, error) {
	inviter, err := app.users.User(inviterId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	invitee, err := app.users.UserByUsername(inviteeUsername)
	if err != nil {
		return invitation.Invitation{}, err
	}

	if inviter.Id == invitee.Id {
		return invitation.Invitation{}, errors.New("inviter cannot be the same as invitee")
	}

	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	inv := invitation.New(
		invitation.WithKeystoreId(k.Id),
		invitation.WithInviterId(inviter.Id),
		invitation.WithInviteeId(invitee.Id),
	)

	return app.invitations.CreateInvitation(inv)
}

func (app *application) AcceptInvitation(userId string, invitationId string) (invitation.Invitation, error) {
	inv, err := app.UserInvitation(userId, invitationId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	if userId != inv.InviteeId {
		return invitation.Invitation{}, errors.New("cannot accept invitation")
	}

	err = inv.Accept()
	if err != nil {
		return invitation.Invitation{}, err
	}

	inv, err = app.invitations.UpdateInvitation(inv)
	if err != nil {
		return invitation.Invitation{}, err
	}

	return inv, nil
}

// func (app *application) VerifyInvitation(userId string, invitationId string) (invitation.Invitation, error) {
// 	u, err := app.users.User(userId)
// 	if err != nil {
// 		return invitation.Invitation{}, err
// 	}
//
// 	inv, err := app.invitations.Invitation(invitationId)
// 	if err != nil {
// 		return invitation.Invitation{}, err
// 	}
//
// 	if u.Id != inv.InviterId && u.Id != inv.InviteeId {
// 		return invitation.Invitation{}, errors.New("unauthorized")
// 	}
//
// 	if u.Id == inv.InviterId {
// 		err = inv.Delete()
// 		if err != nil {
// 			return invitation.Invitation{}, err
// 		}
// 	} else if u.Id == inv.InviteeId {
// 		err = inv.Decline()
// 		if err != nil {
// 			return invitation.Invitation{}, err
// 		}
// 	}
// }

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

	inv, err = app.invitations.UpdateInvitation(inv)
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

	inv, err = app.invitations.UpdateInvitation(inv)
	if err != nil {
		return invitation.Invitation{}, err
	}

	return inv, nil
}
