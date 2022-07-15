package invitation

import (
	"errors"
)

var (
	ErrNotFound = errors.New("invitation not found")
)

type Repository interface {
	CreateInvitation(opts ...Option) (Invitation, error)
	Invitation(id string) (Invitation, error)
	UpdateInvitation(*Invitation) error
	Invitations() ([]Invitation, error)
	Delete(id string) error
	// UserInvitations(userId string) ([]Invitation, error)
	// UserInvitation(userId, invitationId string) (Invitation, error)
}

type UserInvitationsOptions struct {
	Status *Status
}

type Service interface {
	CreateInvitation(keystoreId, inviterId, inviteeUsername string) (Invitation, error)
	AcceptInvitation(userId string, invitationId string) (Invitation, error)
	DeclineOrCancelInvitation(userId, invitationId string) (Invitation, error)
	FinalizeInvitation(invitationId string, encryptedKeystoreKey []byte) (Invitation, error)
	UserInvitation(userId, invitationId string) (Invitation, error)
	UserInvitations(userId string, opts *UserInvitationsOptions) ([]Invitation, error)
}
