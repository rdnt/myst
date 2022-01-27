package keystore

import (
	"myst/internal/client/core/domain/keystore/entry"
)

type Option func(*Keystore)

func WithId(id string) Option {
	return func(k *Keystore) {
		k.id = id
	}
}

func WithName(name string) Option {
	return func(k *Keystore) {
		k.name = name
	}
}

func WithVersion(version int) Option {
	return func(k *Keystore) {
		k.version = version
	}
}

func WithPassword(password string) Option {
	return func(k *Keystore) {
		k.password = password
	}
}

func WithEntries(entries []entry.Entry) Option {
	return func(k *Keystore) {
		k.entries = entries
	}
}
