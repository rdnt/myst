package enclave

import (
	"myst/internal/client/application/domain/keystore"
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

func WithUsername(username string) Option {
	return func(e *Enclave) error {
		e.Username = username

		return nil
	}
}

func WithPassword(password string) Option {
	return func(e *Enclave) error {
		e.Password = password

		return nil
	}
}

func WithPublicKey(publicKey []byte) Option {
	return func(e *Enclave) error {
		e.PublicKey = publicKey

		return nil
	}
}

func WithPrivateKey(privateKey []byte) Option {
	return func(e *Enclave) error {
		e.PrivateKey = privateKey

		return nil
	}
}
