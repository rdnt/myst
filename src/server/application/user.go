package application

import (
	"github.com/pkg/errors"

	"myst/pkg/crypto"
	"myst/src/server/application/domain/invitation"
	"myst/src/server/application/domain/user"
)

var (
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")
)

func (app *application) CreateUser(username, password string, publicKey []byte) (user.User, error) {
	if username == "" {
		return user.User{}, ErrInvalidUsername
	}

	if password == "" {
		return user.User{}, ErrInvalidPassword
	}

	hash, err := crypto.HashPassword(password)
	if err != nil {
		return user.User{}, err
	}

	return app.users.CreateUser(user.New(
		user.WithUsername(username),
		user.WithPasswordHash(hash),
		user.WithPublicKey(publicKey),
	))
}

func (app *application) AuthorizeUser(username, password string) (user.User, error) {
	u, err := app.users.UserByUsername(username)
	if err != nil {
		return user.User{}, err
	}

	ok, err := crypto.VerifyPassword(password, u.PasswordHash)
	if err != nil {
		return user.User{}, err
	}

	if !ok {
		return user.User{}, errors.New("invalid password")
	}

	return u, nil
}

func (app *application) User(userId string) (user.User, error) {
	return app.users.User(userId)
}

func (app *application) UserByUsername(username string) (user.User, error) {
	return app.users.UserByUsername(username)
}

func (app *application) UserInvitations(userId string, opts UserInvitationsOptions) ([]invitation.Invitation, error) {
	u, err := app.users.User(userId)
	if err != nil {
		return nil, err
	}

	invs, err := app.invitations.Invitations()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get user invitations")
	}

	invitations := []invitation.Invitation{}
	for _, inv := range invs {
		if inv.InviterId != u.Id && inv.InviteeId != u.Id {
			continue
		}

		if inv.Deleted() && inv.InviteeId == u.Id {
			continue
		}

		if opts.Status != nil && *opts.Status != inv.Status {
			continue
		}

		invitations = append(invitations, inv)
	}

	return invitations, nil
}

func (app *application) UserInvitation(userId, invitationId string) (invitation.Invitation, error) {
	u, err := app.users.User(userId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	i, err := app.invitations.Invitation(invitationId)
	if err != nil {
		return invitation.Invitation{}, err
	}

	if i.InviterId != u.Id && i.InviteeId != u.Id {
		return invitation.Invitation{}, invitation.ErrNotFound
	}

	if i.Deleted() && i.InviteeId == u.Id {
		return invitation.Invitation{}, invitation.ErrNotFound
	}

	return i, nil
}
