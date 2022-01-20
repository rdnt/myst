package keystore

import (
	"errors"
)

var (
	ErrNotFound = errors.New("keystore not found")
)

type Repository interface {
	Create(opts ...Option) (*Keystore, error)
	Keystore(id string) (*Keystore, error)
	Update(k *Keystore) error
	Keystores() ([]*Keystore, error)
	Delete(id string) error

	UserKeystore(userId, keystoreId string) (*Keystore, error)
	UserKeystores(userId string) ([]*Keystore, error)
}

type Service interface {
	Create(name, ownerId string, payload []byte) (*Keystore, error)
	Keystore(id string) (*Keystore, error)
	Keystores() ([]*Keystore, error)

	UserKeystore(userId, keystoreId string) (*Keystore, error)
	UserKeystores(userId string) ([]*Keystore, error)
}
