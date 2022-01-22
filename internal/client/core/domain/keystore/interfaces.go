package keystore

import "fmt"

var (
	ErrAuthenticationRequired = fmt.Errorf("authentication required")
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
	Keystore(id string) (*Keystore, error)
	Keystores() ([]*Keystore, error)
	Unlock(id string, pass string) (*Keystore, error)
	Update(k *Keystore) error
}
