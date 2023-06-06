package application

import (
	"time"

	"github.com/pkg/errors"

	"myst/src/server/application/domain/invitation"
)

var (
	ErrInvalidInvitee = errors.New("invalid inviter")
)

// CreateInvitation creates an invitation for the given keystoreId, from the
// given inviter to the invitee. The inviter should be the owner of the keystore.
// Errors returned:
// - ErrUserNotFound if the inviter or invitee are not found.
// - ErrKeystoreNotFound if the keystore is not found.
// - ErrNotAllowed if the inviter is not the owner of the keystore.
// - ErrInvalidInvitee if the inviter and invitee are the same user.
func (app *application) CreateInvitation(keystoreId, inviterId, inviteeUsername string) (invitation.Invitation, error) {
	inviter, err := app.users.User(inviterId)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get inviter by id")
	}

	k, err := app.keystores.Keystore(keystoreId)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get keystore by id")
	}

	if k.OwnerId != inviter.Id {
		return invitation.Invitation{}, ErrNotAllowed
	}

	invitee, err := app.users.UserByUsername(inviteeUsername)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get invitee by username")
	}

	if inviter.Id == invitee.Id {
		return invitation.Invitation{}, ErrInvalidInvitee
	}

	inv := invitation.New(
		invitation.WithKeystoreId(k.Id),
		invitation.WithInviterId(inviter.Id),
		invitation.WithInviteeId(invitee.Id),
	)

	inv, err = app.invitations.CreateInvitation(inv)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to create invitation")
	}

	return inv, nil
}

// AcceptInvitation allows a user to accept an invitation. A user can only
// accept an invitation if they are the invitee and the invitation's status is
// pending. If either is not true then ErrNotAllowed is returned.
// ErrInvitationNotFound is returned if the user does not have access
// to the invitation.
func (app *application) AcceptInvitation(userId string, invitationId string) (invitation.Invitation, error) {
	inv, err := app.UserInvitation(userId, invitationId)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get user invitation")
	}

	if userId != inv.InviteeId {
		return invitation.Invitation{}, ErrNotAllowed
	}

	if inv.Status != invitation.Pending {
		return invitation.Invitation{}, ErrNotAllowed
	}

	inv.Status = invitation.Accepted
	inv.AcceptedAt = time.Now()

	inv, err = app.invitations.UpdateInvitation(inv)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to update invitation")
	}

	return inv, nil
}

// DeleteInvitation deletes, declines or cancels an invitation, depending on
// the type of association of the user against the invitation (inviter:
// deletes/cancels, invitee: declines) and the invitation's current status.
// ErrInvitationNotFound will be returned if the user is not associated with
// the invitation.
// ErrNotAllowed is returned if the invitation is in a state where the operation
// cannot be performed. See deleteInvitation, cancelInvitation and
// declineInvitation for more details on the respective conditions.
func (app *application) DeleteInvitation(userId, invitationId string) (invitation.Invitation, error) {
	inv, err := app.invitations.Invitation(invitationId)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get invitation")
	}

	if userId != inv.InviterId && userId != inv.InviteeId {
		return invitation.Invitation{}, ErrInvitationNotFound
	}

	if userId == inv.InviterId {
		inv, err = app.deleteOrCancelInvitation(inv)
		if err != nil {
			return invitation.Invitation{}, errors.WithMessage(err, "failed to delete or cancel invitation")
		}
	} else if userId == inv.InviteeId {
		inv, err = app.declineInvitation(inv)
		if err != nil {
			return invitation.Invitation{}, errors.WithMessage(err, "failed to decline invitation")
		}
	}

	return inv, nil
}

func (app *application) deleteOrCancelInvitation(inv invitation.Invitation) (invitation.Invitation, error) {
	if inv.Pending() {
		inv.Status = invitation.Deleted
		inv.DeletedAt = time.Now()
	} else if inv.Accepted() {
		inv.Status = invitation.Cancelled
		inv.CancelledAt = time.Now()
	} else {
		return invitation.Invitation{}, ErrNotAllowed
	}

	var err error
	inv, err = app.invitations.UpdateInvitation(inv)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to update invitation")
	}

	return inv, nil
}

func (app *application) declineInvitation(inv invitation.Invitation) (invitation.Invitation, error) {
	if !inv.Pending() {
		return invitation.Invitation{}, ErrNotAllowed
	}

	inv.Status = invitation.Declined
	inv.DeclinedAt = time.Now()

	var err error
	inv, err = app.invitations.UpdateInvitation(inv)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to update invitation")
	}

	return inv, nil
}

// FinalizeInvitation finalizes an invitation. setting the encrypted keystore key
// for the keystore associated with the invitation, allowing decryption of its
// payload by the invitee. If the user doesn't have access to the invitation,
// ErrInvitationNotFound is returned. ErrNotAllowed is returned if the user is
// not the inviter or the invitation is not in the accepted state.
func (app *application) FinalizeInvitation(userId, invitationId string, encryptedKeystoreKey []byte) (invitation.Invitation, error) {
	inv, err := app.invitations.Invitation(invitationId)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get invitation")
	}

	if userId != inv.InviterId && userId != inv.InviteeId {
		return invitation.Invitation{}, ErrInvitationNotFound
	}

	if userId != inv.InviterId {
		return invitation.Invitation{}, ErrNotAllowed
	}

	if !inv.Accepted() {
		return invitation.Invitation{}, ErrNotAllowed
	}

	inv.EncryptedKeystoreKey = encryptedKeystoreKey
	inv.Status = invitation.Finalized
	inv.FinalizedAt = time.Now()

	inv, err = app.invitations.UpdateInvitation(inv)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to update invitation")
	}

	return inv, nil
}
