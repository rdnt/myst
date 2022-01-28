package keystore

type Repository interface {
	Create(opts ...Option) (*Keystore, error)
	Keystore(id string) (*Keystore, error)
	Update(k *Keystore) error
	Keystores() (map[string]*Keystore, error)
	Delete(id string) error
}

type Service interface {
	Create(name string) (*Keystore, error)
	Initialize(name, password string) (*Keystore, error)
	Keystore(id string) (*Keystore, error)
	Keystores() (map[string]*Keystore, error)
	Update(k *Keystore) error
	Authenticate(password string) error
	HealthCheck()
}
