package application

import (
	"myst/internal/server/core/domain/invitation"
	"myst/internal/server/core/domain/keystore"

	"github.com/pkg/errors"
)

func (app *application) CreateInvitation(keystoreId, inviterId, inviteeId string, inviterKey []byte) (*invitation.Invitation, error) {
	return app.Invitations.Create(keystoreId, inviterId, inviteeId, inviterKey)
}

func (app *application) AcceptInvitation(invitationId string, inviteeKey []byte) (*invitation.Invitation, error) {
	return app.Invitations.Accept(invitationId, inviteeKey)
}

func (app *application) FinalizeInvitation(invitationId string, keystoreKey []byte) (*invitation.Invitation, error) {
	return app.Invitations.Finalize(invitationId, keystoreKey)
}

func (app *application) GetInvitation(invitationId string) (*invitation.Invitation, error) {
	return app.Invitations.Invitation(invitationId)
}

func (app *application) CreateKeystore(name, ownerId string, payload []byte) (*keystore.Keystore, error) {
	return app.Keystores.Create(name, ownerId, payload)
}

func (app *application) UserKeystores(userId string) ([]*keystore.Keystore, error) {
	keystores, err := app.Keystores.UserKeystores(userId)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get user keystores")
	}

	status := invitation.Finalized
	invs, err := app.Invitations.UserInvitations(userId, &invitation.UserInvitationsOptions{Status: &status})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get user invitations")
	}

	for _, inv := range invs {
		k, err := app.Keystores.Keystore(inv.KeystoreId())
		if err != nil {
			return nil, errors.WithMessage(err, "failed to get keystore")
		}
		keystores = append(keystores, k)
	}

	return keystores, nil
}

func (app *application) UserInvitations(userId string) ([]*invitation.Invitation, error) {
	invs, err := app.Invitations.UserInvitations(userId, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get user invitations")
	}

	return invs, nil
}

func (app *application) UserKeystore(userId, keystoreId string) (*keystore.Keystore, error) {
	k, err := app.Keystores.UserKeystore(userId, keystoreId)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get user keystore")
	}

	return k, nil
}
