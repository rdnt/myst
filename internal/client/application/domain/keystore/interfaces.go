package keystore

import (
	"myst/internal/client/application/domain/entry"
)

type Repository interface {
	CreateKeystore(k Keystore) (Keystore, error)
	Keystore(id string) (Keystore, error)
	UpdateKeystore(k Keystore) error
	Keystores() (map[string]Keystore, error)
	DeleteKeystore(id string) error

	Authenticate(password string) error
	CreateEnclave(password string) error
	HealthCheck()
	//KeystoreKey(keystoreId string) ([]byte, error)
}

type Service interface {
	CreateKeystore(k Keystore) (Keystore, error)
	DeleteKeystore(id string) error
	CreateFirstKeystore(k Keystore, password string) (Keystore, error)
	Keystore(id string) (Keystore, error)
	KeystoreByRemoteId(id string) (Keystore, error)
	KeystoreEntries(id string) (map[string]entry.Entry, error)
	CreateKeystoreEntry(keystoreId string, opts ...entry.Option) (entry.Entry, error)
	UpdateKeystoreEntry(keystoreId string, entryId string, password, notes *string) (entry.Entry, error)
	DeleteKeystoreEntry(keystoreId, entryId string) error
	Keystores() (map[string]Keystore, error)
	UpdateKeystore(k Keystore) error
	Authenticate(password string) error
	HealthCheck()
	Keypair() (publicKey, privateKey []byte, err error)
	SetKeypair(publicKey, privateKey []byte) error
	UserInfo() (username, password string, err error)
	SetUserInfo(username, password string) error
	CreateEnclave(password string) error
	Enclave() error
}
