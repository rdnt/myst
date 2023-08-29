package application

import (
	"myst/src/client/application/domain/credentials"
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/user"
)

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

type Remote interface {
	Address() string

	Authenticate(username, password string) error
	Register(username, password string, publicKey []byte) (user.User, error)
	Authenticated() bool
	CurrentUser() *user.User

	CreateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	UpdateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	Keystores(privateKey []byte) (map[string]keystore.Keystore, error)
	DeleteKeystore(id string) error

	CreateInvitation(keystoreRemoteId, inviteeUsername string) (invitation.Invitation, error)
	Invitation(id string) (invitation.Invitation, error)
	AcceptInvitation(id string) (invitation.Invitation, error)
	DeleteInvitation(id string) (invitation.Invitation, error)
	FinalizeInvitation(invitationId string, privateKey, inviteePublicKey []byte,
		k keystore.Keystore) (invitation.Invitation, error)
	Invitations() (map[string]invitation.Invitation, error)

	SharedSecret(privateKey []byte, userId string) ([]byte, error)
}
