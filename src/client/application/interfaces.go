package application

import (
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/user"
)

// Remote is a remote repository that holds upstream keystores/invitations
type Remote interface {
	Address() string
	CreateInvitation(inv invitation.Invitation) (invitation.Invitation, error)
	Invitation(id string) (invitation.Invitation, error)
	AcceptInvitation(id string) (invitation.Invitation, error)
	DeclineOrCancelInvitation(id string) (invitation.Invitation, error)
	FinalizeInvitation(invitationId string, keystoreKey, privateKey []byte) (invitation.Invitation, error)
	UpdateInvitation(inv invitation.Invitation) error
	Invitations() (map[string]invitation.Invitation, error)
	DeleteInvitation(id string) error

	CreateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	Keystore(id string, privateKey []byte) (keystore.Keystore, error)
	UpdateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	Keystores(privateKey []byte) (map[string]keystore.Keystore, error)
	DeleteKeystore(id string) error

	SignIn(username, password string, publicKey []byte) (user.User, error)
	Register(username, password string, publicKey []byte) (user.User, error)
	SignOut() error
	SignedIn() bool
	CurrentUser() *user.User
	UserByUsername(username string) (user.User, error)
}

type Repository interface {
	HealthCheck()

	Initialize(password string) error
	IsInitialized() error

	Authenticate(password string) error
}
