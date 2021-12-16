package server

import (
	"myst/app/server/core/domain/invitation"
)

func (app *Application) CreateKeystoreInvitation(
	inviterId, inviteeId, keystoreId string, inviterKey []byte,
) (*invitation.Invitation, error) {
	inviter, err := app.userRepo.User(inviterId)
	if err != nil {
		return nil, err
	}

	invitee, err := app.userRepo.User(inviteeId)
	if err != nil {
		return nil, err
	}

	ks, err := app.keystoreRepo.Keystore(keystoreId)
	if err != nil {
		return nil, err
	}

	inv, err := app.invitationService.Create(
		invitation.WithKeystore(ks),
		invitation.WithInviter(inviter),
		invitation.WithInvitee(invitee),
		invitation.WithInviterKey(inviterKey),
	)
	if err != nil {
		return nil, err
	}

	return inv, nil
}
