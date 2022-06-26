package keystore

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
