package server

import (
	"encoding/hex"

	"myst/app/server/domain/invitation"
)

func (app *Application) CreateKeystoreInvitation(
	inviterId, inviteeId, keystoreId, inviterKey string,
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

	key, err := hex.DecodeString(inviterKey)
	if err != nil {
		return nil, err
	}

	inv, err := app.invitationService.Create(
		invitation.WithKeystore(ks),
		invitation.WithInviter(inviter),
		invitation.WithInvitee(invitee),
		invitation.WithInviterKey(key),
	)
	if err != nil {
		return nil, err
	}

	return inv, nil
}
