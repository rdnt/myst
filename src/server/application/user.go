package application

import (
	"github.com/pkg/errors"

	"myst/src/server/application/domain/invitation"
	"myst/src/server/application/domain/user"
)

func (app *application) CreateUser(username, password string, publicKey []byte) (user.User, error) {
	return app.users.CreateUser(
		user.WithUsername(username),
		user.WithPassword(password),
		user.WithPublicKey(publicKey),
	)
}

func (app *application) AuthorizeUser(username, password string) error {
	u, err := app.users.UserByUsername(username)
	if err != nil {
		return err
	}

	ok, err := app.users.VerifyPassword(u.Id, password)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("invalid password")
	}

	return nil
}

func (app *application) User(userId string) (user.User, error) {
	return app.users.User(userId)
}

func (app *application) UserByUsername(username string) (user.User, error) {
	return app.users.UserByUsername(username)
}

func (app *application) UserInvitations(userId string, opts *invitation.UserInvitationsOptions) ([]invitation.Invitation, error) {
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
		if inv.Deleted() && inv.InviteeId == u.Id {
			continue
		}
		if opts != nil && opts.Status != nil && *opts.Status != inv.Status {
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

	return i, nil
}
