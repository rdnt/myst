package application

import (
	"myst/internal/server/core/domain/invitation"
	"myst/internal/server/core/domain/keystore"
)

func (app *Application) CreateInvitation(
	keystoreId, inviterId, inviteeId string, inviterKey []byte,
) (*invitation.Invitation, error) {
	return app.Invitations.Create(keystoreId, inviterId, inviteeId, inviterKey)
}

func (app *Application) AcceptInvitation(invitationId string, inviteeKey []byte) (*invitation.Invitation, error) {
	return app.Invitations.Accept(invitationId, inviteeKey)
}

func (app *Application) FinalizeInvitation(invitationId string, keystoreKey []byte) (*invitation.Invitation, error) {
	return app.Invitations.Finalize(invitationId, keystoreKey)
}

func (app *Application) GetInvitation(invitationId string) (*invitation.Invitation, error) {
	return app.Invitations.Invitation(invitationId)
}

func (app *Application) CreateKeystore(name, ownerId string, payload []byte) (*keystore.Keystore, error) {
	return app.Keystores.Create(name, ownerId, payload)
}
