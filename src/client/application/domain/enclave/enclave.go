package enclave

import (
	"fmt"

	"myst/pkg/crypto"
	"myst/src/client/application/domain/keystore"
)

var (
	ErrRemoteNotSet = fmt.Errorf("remote not set")
)

type Enclave struct {
	salt      []byte
	keystores map[string]keystore.Keystore
	keys      map[string][]byte

	remote *Remote
}

type Remote struct {
	Address    string
	Username   string
	Password   string
	PublicKey  []byte
	PrivateKey []byte
}

func New(opts ...Option) (*Enclave, error) {
	e := &Enclave{
		keystores: map[string]keystore.Keystore{},
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

func (e *Enclave) AddKeystore(k keystore.Keystore) error {
	p := crypto.DefaultArgon2IdParams

	key, err := crypto.GenerateRandomBytes(uint(p.KeyLength))
	if err != nil {
		return err
	}

	e.keystores[k.Id] = k
	e.keys[k.Id] = key

	return nil
}

func (e *Enclave) Keys() map[string][]byte {
	return e.keys
}

func (e *Enclave) Keystores() (map[string]keystore.Keystore, error) {
	ks := map[string]keystore.Keystore{}

	for _, k := range e.keystores {
		keystoreKey, ok := e.keys[k.Id]
		if !ok {
			return nil, fmt.Errorf("keystore key not found")
		}

		k.Key = keystoreKey
		ks[k.Id] = k
	}

	return ks, nil
}

func (e *Enclave) Keystore(id string) (keystore.Keystore, error) {
	k, ok := e.keystores[id]
	if !ok {
		return keystore.Keystore{}, fmt.Errorf("keystore not found")
	}

	keystoreKey, ok := e.keys[id]
	if !ok {
		return keystore.Keystore{}, fmt.Errorf("keystore key not found")
	}

	k.Key = keystoreKey

	return k, nil
}

func (e *Enclave) Salt() []byte {
	return e.salt
}

func (e *Enclave) UpdateKeystore(k keystore.Keystore) error {
	e.keystores[k.Id] = k

	return nil
}

func (e *Enclave) DeleteKeystore(id string) error {
	delete(e.keystores, id)
	delete(e.keys, id)

	return nil
}

func (e *Enclave) SetRemote(r Remote) {
	e.remote = &r
}

func (e *Enclave) Remote() *Remote {
	return e.remote
}