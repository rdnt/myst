package application

import (
	"myst/src/client/application/domain/credentials"
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/user"
)

// Enclave is the repository that handles storing and retrieving of the
// user's keystores and credentials. It requires initialization and
// authentication before it can be used. The authentication status can
// expire after some time if the HealthCheck method is not called
// regularly.
type Enclave interface {
	Initialize(password string) (key []byte, err error)
	IsInitialized() (bool, error)
	Authenticate(password string) (key []byte, err error)

	CreateKeystore(key []byte, k keystore.Keystore) (keystore.Keystore, error)
	Keystore(key []byte, id string) (keystore.Keystore, error)
	UpdateKeystore(key []byte, k keystore.Keystore) (keystore.Keystore, error)
	Keystores(key []byte) (map[string]keystore.Keystore, error)
	DeleteKeystore(key []byte, id string) error

	UpdateCredentials(key []byte, creds credentials.Credentials) (credentials.Credentials, error)
	Credentials(key []byte) (credentials.Credentials, error)
}

// Remote is a remote repository that holds upstream enclave/invitations. It is
// used to sync keystores with a remote server in a secure manner, and to
// facilitate inviting users to access keystores or accepting invitations to
// access a keystore from another user. Authenticating with a username and
// password is required to interface with a remote.
type Remote interface {
	Address() string

	CreateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	UpdateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	Keystores(privateKey []byte) (map[string]keystore.Keystore, error)
	DeleteKeystore(id string) error

	CreateInvitation(keystoreRemoteId, inviteeUsername string) (invitation.Invitation, error)
	Invitation(id string) (invitation.Invitation, error)
	AcceptInvitation(id string) (invitation.Invitation, error)
	DeleteInvitation(id string) (invitation.Invitation, error)
	FinalizeInvitation(invitationId string, privateKey, inviteePublicKey []byte, k keystore.Keystore) (invitation.Invitation, error)
	Invitations() (map[string]invitation.Invitation, error)

	Authenticate(username, password string) error
	Register(username, password string, publicKey []byte) (user.User, error)
	Authenticated() bool
	CurrentUser() *user.User

	SharedSecret(privateKey []byte, userId string) ([]byte, error)
}
