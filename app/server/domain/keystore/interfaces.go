package keystore

type Repository interface {
	CreateKeystore(opts ...Option) (*Keystore, error)
	Keystore(id string) (*Keystore, error)
	UpdateKeystore(*Keystore) error
	Keystores() ([]*Keystore, error)
	DeleteKeystore(id string) error
}
