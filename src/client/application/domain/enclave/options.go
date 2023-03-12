package enclave

import (
	"myst/src/client/application/domain/credentials"
	"myst/src/client/application/domain/keystore"
)

type Option func(*Enclave) error

func WithKeystores(keystores map[string]keystore.Keystore) Option {
	return func(e *Enclave) error {
		for _, k := range keystores {
			err := e.AddKeystore(k)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func WithSalt(salt []byte) Option {
	return func(e *Enclave) error {
		e.salt = salt

		return nil
	}
}

func WithRemote(remote *credentials.Credentials) Option {
	return func(e *Enclave) error {
		e.remote = remote

		return nil
	}
}
