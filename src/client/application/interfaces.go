package application

import (
	"myst/src/client/application/domain/credentials"
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/user"
)

// Enclave is the repository that handles storing and retrieving of the
// user's keystores and credentials. It requires initialization and
// authentication before it can be used. Authentication status can
// expire after some time if the HealthCheck method is not called
// regularly.
type Enclave interface {
	Initialize(password string) error
	IsInitialized() (bool, error)
	Authenticate(password string) error
	HealthCheck()

	CreateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	Keystore(id string) (keystore.Keystore, error)
	UpdateKeystore(k keystore.Keystore) error
	Keystores() (map[string]keystore.Keystore, error)
	DeleteKeystore(id string) error

	SetCredentials(address, username, password string, publicKey, privateKey []byte) error
	Credentials() (credentials.Credentials, error)
}

// Remote is a remote repository that holds upstream enclave/invitations. It is
// used to sync keystores with a remote server in a secure manner, and to
// facilitate inviting users to access keystores or accepting invitations to
// access a keystore from another user. Authentication with a username and
// password is required to interface with a remote.
type Remote interface {
	Address() string

	CreateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	Keystore(id string, privateKey []byte) (keystore.Keystore, error)
	UpdateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	Keystores(privateKey []byte) (map[string]keystore.Keystore, error)
	DeleteKeystore(id string) error

	CreateInvitation(inv invitation.Invitation) (invitation.Invitation, error)
	Invitation(id string) (invitation.Invitation, error)
	AcceptInvitation(id string) (invitation.Invitation, error)
	DeclineOrCancelInvitation(id string) (invitation.Invitation, error)
	FinalizeInvitation(invitationId string, encryptedKeystoreKey []byte) (invitation.Invitation, error)
	Invitations() (map[string]invitation.Invitation, error)

	SignIn(username, password string) error
	Register(username, password string, publicKey []byte) (user.User, error)
	// SignOut() error
	SignedIn() bool
	CurrentUser() *user.User
	UserByUsername(username string) (user.User, error)
}
