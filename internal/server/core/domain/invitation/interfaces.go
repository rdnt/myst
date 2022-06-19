package invitation

import (
	"errors"
)

var (
	ErrNotFound = errors.New("invitation not found")
)

type Repository interface {
	Create(opts ...Option) (*Invitation, error)
	Invitation(id string) (*Invitation, error)
	Update(k *Invitation) error
	Invitations() ([]*Invitation, error)
	Delete(id string) error

	UserInvitations(userId string) ([]Invitation, error)
	UserInvitation(userId, invitationId string) (*Invitation, error)
}

type UserInvitationsOptions struct {
	Status *Status
}

type Service interface {
	Create(keystoreId, inviterId, inviteeId string, inviterKey []byte) (*Invitation, error)
	Invitation(id string) (*Invitation, error)

	// TODO: pass user id with them to verify the user exists
	Accept(invitationId string, inviteeKey []byte) (*Invitation, error)
	Finalize(invitationId string, keystoreKey []byte) (*Invitation, error)

	UserInvitations(userId string, opts *UserInvitationsOptions) ([]Invitation, error)
	UserInvitation(userId, invitationId string) (*Invitation, error)
	DeclineOrCancelInvitation(userId, invitationId string) (*Invitation, error)

	//UserKeystores(userId string) ([]*keystore.Keystore, error)
	//AcceptedUserInvitations(userId string) ([]*Invitation, error)
}
