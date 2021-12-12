package keystore

import (
	"myst/app/client/core/domain/keystore/entry"
)

type Option func(*Keystore) error

func WithVersion(version int) Option {
	return func(k *Keystore) error {
		k.version = version
		return nil
	}
}

func WithEntries(entries []*entry.Entry) Option {
	return func(k *Keystore) error {
		k.entries = entries
		return nil
	}
}
