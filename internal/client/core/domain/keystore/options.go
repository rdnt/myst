package keystore

import (
	"myst/internal/client/core/domain/keystore/entry"
)

type Option func(*Keystore) error

func WithId(id string) Option {
	return func(k *Keystore) error {
		k.id = id
		return nil
	}
}

func WithName(name string) Option {
	return func(k *Keystore) error {
		k.name = name
		return nil
	}
}

func WithVersion(version int) Option {
	return func(k *Keystore) error {
		k.version = version
		return nil
	}
}

func WithPassphrase(passphrase string) Option {
	return func(k *Keystore) error {
		k.passphrase = passphrase
		return nil
	}
}

func WithEntries(entries []*entry.Entry) Option {
	return func(k *Keystore) error {
		k.entries = entries
		return nil
	}
}
