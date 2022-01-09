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
}

type Service interface {
	Create(opts ...Option) (*Keystore, error)
}
