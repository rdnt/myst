package keystore

import "myst/internal/client/application/domain/entry"

type Repository interface {
	Create(opts ...Option) (*Keystore, error)
	Keystore(id string) (*Keystore, error)
	Update(k *Keystore) error
	Keystores() (map[string]*Keystore, error)
	Delete(id string) error
}

type Service interface {
	CreateKeystore(name string) (*Keystore, error)
	CreateFirstKeystore(name, password string) (*Keystore, error)
	Keystore(id string) (*Keystore, error)
	KeystoreEntries(id string) (map[string]entry.Entry, error)
	CreateKeystoreEntry(keystoreId string, opts ...entry.Option) (entry.Entry, error)
	UpdateKeystoreEntry(keystoreId string, entryId string, password, notes *string) (entry.Entry, error)
	DeleteKeystoreEntry(keystoreId, entryId string) error
	Keystores() (map[string]*Keystore, error)
	UpdateKeystore(k *Keystore) error
	Authenticate(password string) error
	HealthCheck()
}
