package enclaverepo

import (
	"fmt"
	"github.com/pkg/errors"
	"myst/src/server/application"
	"time"

	"myst/pkg/crypto"
	"myst/src/client/application/domain/credentials"
	"myst/src/client/application/domain/keystore"
)

var (
	// TODO @rdnt @@@: move to app if possible
	ErrRemoteNotSet = fmt.Errorf("remote not set")
)

type Enclave struct {
	salt      []byte
	keystores map[string]keystore.Keystore
	keys      map[string][]byte

	remote *credentials.Credentials
}

func newEnclave(opts ...Option) *Enclave {
	e := &Enclave{
		keystores: map[string]keystore.Keystore{},
		keys:      map[string][]byte{},
	}

	for _, opt := range opts {
		if opt != nil {
			opt(e)
		}
	}

	return e
}

type Option func(*Enclave)

func WithKeystores(keystores map[string]keystore.Keystore) Option {
	return func(e *Enclave) {
		e.keystores = keystores
	}
}

func WithKeys(keys map[string][]byte) Option {
	return func(e *Enclave) {
		e.keys = keys
	}
}

func WithSalt(salt []byte) Option {
	return func(e *Enclave) {
		e.salt = salt
	}
}

func WithRemote(remote *credentials.Credentials) Option {
	return func(e *Enclave) {
		e.remote = remote
	}
}

func (e *Enclave) AddKeystore(k keystore.Keystore) error {
	p := crypto.DefaultArgon2IdParams

	key, err := crypto.GenerateRandomBytes(uint(p.KeyLength))
	if err != nil {
		return errors.WithMessage(err, "failed to generate key")
	}

	k.CreatedAt = time.Now()
	k.UpdatedAt = time.Now()

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
			return nil, application.ErrKeystoreNotFound
		}

		k.Key = keystoreKey
		ks[k.Id] = k
	}

	return ks, nil
}

func (e *Enclave) Keystore(id string) (keystore.Keystore, error) {
	k, ok := e.keystores[id]
	if !ok {
		return keystore.Keystore{}, application.ErrKeystoreNotFound
	}

	keystoreKey, ok := e.keys[id]
	if !ok {
		return keystore.Keystore{}, application.ErrKeystoreNotFound
	}

	k.Key = keystoreKey

	return k, nil
}

func (e *Enclave) Salt() []byte {
	return e.salt
}

func (e *Enclave) UpdateKeystore(k keystore.Keystore) error {
	k.UpdatedAt = time.Now()
	k.Version++

	e.keystores[k.Id] = k

	return nil
}

func (e *Enclave) DeleteKeystore(id string) error {
	delete(e.keystores, id)
	delete(e.keys, id)

	return nil
}

func (e *Enclave) SetRemote(r credentials.Credentials) {
	e.remote = &r
}

func (e *Enclave) Remote() *credentials.Credentials {
	return e.remote
}
