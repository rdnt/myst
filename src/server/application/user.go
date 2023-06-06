package application

import (
	"github.com/pkg/errors"

	"myst/pkg/crypto"
	"myst/src/server/application/domain/invitation"
	"myst/src/server/application/domain/user"
)

// CreateUser creates a user with the given username, password and public key.
// If the username or password are invalid, ErrInvalidUsername or
// ErrInvalidPassword will be returned respectively.
func (app *application) CreateUser(username, password string, publicKey []byte) (user.User, error) {
	if username == "" {
		return user.User{}, ErrInvalidUsername
	}

	if password == "" {
		return user.User{}, ErrInvalidPassword
	}

	hash, err := crypto.HashPassword(password)
	if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to hash password")
	}

	u, err := app.users.CreateUser(user.New(
		user.WithUsername(username),
		user.WithPasswordHash(hash),
		user.WithPublicKey(publicKey),
	))
	if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to create user")
	}

	return u, nil
}

// AuthorizeUser authorizes a user with the given username and password.
// If the password does not match, ErrInvalidPassword is returned.
func (app *application) AuthorizeUser(username, password string) (user.User, error) {
	u, err := app.users.UserByUsername(username)
	if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to get user by username")
	}

	ok, err := crypto.VerifyPassword(password, u.PasswordHash)
	if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to verify password")
	}

	if !ok {
		return user.User{}, ErrInvalidPassword
	}

	return u, nil
}

// User returns a user.
func (app *application) User(userId string) (user.User, error) {
	u, err := app.users.User(userId)
	if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to get user")
	}

	return u, nil
}

// UserByUsername returns a user by username.
func (app *application) UserByUsername(username string) (user.User, error) {
	u, err := app.users.UserByUsername(username)
	if err != nil {
		return user.User{}, errors.WithMessage(err, "failed to get user by username")
	}

	return u, nil
}

// UserInvitations returns a list of invitations for the given user.
// UserInvitationsOptions can be passed to filter out invitations by status.
// If the invitation is deleted, it is only included in the response if the user
// in question is the inviter.
func (app *application) UserInvitations(userId string, opts UserInvitationsOptions) ([]invitation.Invitation, error) {
	u, err := app.users.User(userId)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get user")
	}

	invs, err := app.invitations.Invitations()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get invitations")
	}

	invitations := []invitation.Invitation{}
	for _, inv := range invs {
		// skip invitations irrelevant to the user
		if inv.InviterId != u.Id && inv.InviteeId != u.Id {
			continue
		}

		// do not show deleted invitations if user is the invitee
		if inv.Deleted() && inv.InviteeId == u.Id {
			continue
		}

		// skip invitations not matching the status filter, if provided
		if opts.Status != nil && *opts.Status != inv.Status {
			continue
		}

		invitations = append(invitations, inv)
	}

	return invitations, nil
}

// UserInvitation returns an invitation, if the user is associated with it.
// If the invitation is not found, ErrInvitationNotFound is returned.
// Additionally, if the invitation is deleted and the user is the invitee,
// ErrInvitationNotFound is returned.
func (app *application) UserInvitation(userId, invitationId string) (invitation.Invitation, error) {
	u, err := app.users.User(userId)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get user")
	}

	i, err := app.invitations.Invitation(invitationId)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "failed to get invitation")
	}

	// do not show invitations irrelevant to the user
	if i.InviterId != u.Id && i.InviteeId != u.Id {
		return invitation.Invitation{}, ErrInvitationNotFound
	}

	// do not show deleted invitations if user is the invitee
	if i.Deleted() && i.InviteeId == u.Id {
		return invitation.Invitation{}, ErrInvitationNotFound
	}

	return i, nil
}
