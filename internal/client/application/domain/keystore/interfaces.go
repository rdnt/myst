package keystore

import "fmt"

var (
	ErrAuthenticationRequired = fmt.Errorf("authentication required")
)

type Repository interface {
	Create(opts ...Option) (*Keystore, error)
	Keystore(id string) (*Keystore, error)
	Update(k *Keystore) error
	KeystoreIds() ([]string, error)
	Keystores() (map[string]*Keystore, error)
	Delete(id string) error
}

type Service interface {
	Create(name string) (*Keystore, error)
	Initialize(name, password string) (*Keystore, error)
	Keystore(id string) (*Keystore, error)
	KeystoreIds() ([]string, error)
	Keystores() (map[string]*Keystore, error)
	//Unlock(id string, pass string) (*Keystore, error)
	Update(k *Keystore) error
	Authenticate(password string) error
}
