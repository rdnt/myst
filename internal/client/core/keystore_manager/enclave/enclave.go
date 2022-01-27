package enclave

import (
	"fmt"

	"myst/internal/client/core/domain/keystore"
	"myst/pkg/crypto"
)

type Enclave struct {
	salt      []byte
	keystores map[string]*keystore.Keystore
	keys      map[string][]byte
}

func New(opts ...Option) (*Enclave, error) {
	e := &Enclave{
		keystores: map[string]*keystore.Keystore{},
		keys:      map[string][]byte{},
	}

	for _, opt := range opts {
		err := opt(e)
		if err != nil {
			return nil, err
		}
	}

	return e, nil
}

func (e *Enclave) AddKeystore(k *keystore.Keystore) error {
	p := crypto.DefaultArgon2IdParams

	key, err := crypto.GenerateRandomBytes(uint(p.KeyLength))
	if err != nil {
		return err
	}

	e.keystores[k.Id()] = k
	e.keys[k.Id()] = key

	return nil
}

func (e *Enclave) Keys() map[string][]byte {
	return e.keys
}

func (e *Enclave) Keystores() map[string]*keystore.Keystore {
	return e.keystores
}

func (e *Enclave) KeystoreIds() []string {
	ids := []string{}

	for _, id := range e.keystores {
		ids = append(ids, id.Id())
	}

	return ids
}

func (e *Enclave) Keystore(id string) (*keystore.Keystore, error) {
	k, ok := e.keystores[id]
	if !ok {
		return nil, fmt.Errorf("keystore not found")
	}

	return k, nil
}

func (e *Enclave) Salt() []byte {
	return e.salt
}

func (e *Enclave) UpdateKeystore(k *keystore.Keystore) error {
	e.keystores[k.Id()] = k

	return nil
}

func (e *Enclave) DeleteKeystore(id string) error {
	delete(e.keystores, id)
	delete(e.keys, id)

	return nil
}
