package enclave

import (
	"myst/internal/client/application/domain/entry"
	"myst/internal/client/application/domain/keystore"
)

type KeystoreService interface {
	CreateKeystore(k keystore.Keystore) (keystore.Keystore, error)
	DeleteKeystore(id string) error
	CreateFirstKeystore(k keystore.Keystore, password string) (keystore.Keystore, error)
	Keystore(id string) (keystore.Keystore, error)
	KeystoreByRemoteId(id string) (keystore.Keystore, error)
	KeystoreEntries(id string) (map[string]entry.Entry, error)
	CreateKeystoreEntry(keystoreId string, opts ...entry.Option) (entry.Entry, error)
	UpdateKeystoreEntry(keystoreId string, entryId string, password, notes *string) (entry.Entry, error)
	DeleteKeystoreEntry(keystoreId, entryId string) error
	Keystores() (map[string]keystore.Keystore, error)
	UpdateKeystore(k keystore.Keystore) error
	Authenticate(password string) error
	HealthCheck()
	CreateEnclave(password string) error
	Enclave() error
	SetRemote(address, username, password string) error
	Remote() (Remote, error)
}
