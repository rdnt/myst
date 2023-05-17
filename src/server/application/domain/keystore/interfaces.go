package keystore

import (
	"errors"
)

var (
	ErrNotFound = errors.New("keystore not found")
)

// type Repository interface {
// 	CreateKeystore(opts ...Option) (Keystore, error)
// 	Keystore(id string) (Keystore, error)
// 	UpdateKeystore(k *Keystore) error
// 	Keystores() ([]Keystore, error)
// 	DeleteKeystore(id string) error
// 	// UserKeystore(userId, keystoreId string) (Keystore, error)
// 	// UserKeystores(userId string) ([]Keystore, error)
// }
//
// type UpdateParams struct {
// 	Name    *string
// 	Payload *[]byte
// }

// type Service interface {
// 	CreateKeystore(name, ownerId string, payload []byte) (Keystore, error)
// 	Keystore(id string) (Keystore, error)
// 	Keystores() ([]Keystore, error)
// 	UpdateKeystore(userId, keystoreId string, params UpdateParams) (Keystore, error)
// 	UserKeystore(userId, keystoreId string) (Keystore, error)
// 	UserKeystores(userId string) ([]Keystore, error)
// }
