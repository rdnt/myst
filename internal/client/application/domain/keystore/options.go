package keystore

import (
	"myst/internal/client/application/domain/entry"
)

type Option func(*Keystore)

func WithId(id string) Option {
	return func(k *Keystore) {
		k.Id = id
	}
}

func WithName(name string) Option {
	return func(k *Keystore) {
		k.Name = name
	}
}

func WithVersion(version int) Option {
	return func(k *Keystore) {
		k.Version = version
	}
}

func WithEntries(entries map[string]entry.Entry) Option {
	return func(k *Keystore) {
		k.Entries = entries
	}
}
