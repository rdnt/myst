package keystore

import (
	"myst/src/client/application/domain/keystore/entry"
)

type Repository interface {
	CreateKeystore(k Keystore) (Keystore, error)
	Keystore(id string) (Keystore, error)
	UpdateKeystore(k Keystore) error
	Keystores() (map[string]Keystore, error)
	DeleteKeystore(id string) error

	// EncryptedKeystoreKey(keystoreId string) ([]byte, error)
}

type Service interface {
	CreateKeystore(k Keystore) (Keystore, error)
	DeleteKeystore(id string) error
	Keystore(id string) (Keystore, error)
	KeystoreByRemoteId(id string) (Keystore, error)
	CreateKeystoreEntry(keystoreId string, opts ...entry.Option) (entry.Entry, error)
	UpdateKeystoreEntry(keystoreId string, entryId string, password, notes *string) (entry.Entry, error)
	DeleteKeystoreEntry(keystoreId, entryId string) error
	Keystores() (map[string]Keystore, error)
	UpdateKeystore(k Keystore) error
}
