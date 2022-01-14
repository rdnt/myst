package application

import (
	"myst/internal/server/core/domain/invitation"
)

func (app *Application) CreateKeystoreInvitation(
	keystoreId, inviterId, inviteeId string, inviterKey []byte,
) (*invitation.Invitation, error) {
	return app.Invitations.Create(keystoreId, inviterId, inviteeId, inviterKey)
}

func (app *Application) AcceptKeystoreInvitation(
	invitationId string, inviteeKey []byte,
) (*invitation.Invitation, error) {
	return app.Invitations.Accept(invitationId, inviteeKey)
}

func (app *Application) FinalizeKeystoreInvitation(
	invitationId string, keystoreKey []byte,
) (*invitation.Invitation, error) {
	return app.Invitations.Finalize(invitationId, keystoreKey)
}
