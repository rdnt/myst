package keystore

import (
	"errors"

	"myst/internal/client/core/domain/keystore/entry"

	"myst/pkg/logger"
	"myst/pkg/uuid"
)

var (
	ErrNotFound      = errors.New("keystore not found")
	ErrEntryNotFound = errors.New("entry not found")
	ErrEntryExists   = errors.New("entry already exists")
)

type Keystore struct {
	id         string
	name       string
	version    int
	entries    []entry.Entry
	passphrase string
}

func (k *Keystore) Id() string {
	return k.id
}

func (k *Keystore) Name() string {
	return k.name
}

func (k *Keystore) SetName(name string) {
	k.name = name
}

func (k *Keystore) Version() int {
	return k.version
}

func (k *Keystore) IncrementVersion() {
	k.version++
}

func (k *Keystore) Entries() []entry.Entry {
	return k.entries
}

func (k *Keystore) AddEntry(entry entry.Entry) error {
	for _, e := range k.entries {
		if e.Id() == entry.Id() {
			return ErrEntryExists
		}
	}

	k.entries = append(k.entries, entry)
	return nil
}

func (k *Keystore) RemoveEntry(id string) error {
	for i, e := range k.entries {
		if e.Id() == id {
			k.entries = append(k.entries[:i], k.entries[i+1:]...)
			return nil
		}
	}

	return ErrEntryNotFound
}

func (k *Keystore) Passphrase() string {
	return k.passphrase
}

func (k *Keystore) SetPassphrase(passphrase string) {
	k.passphrase = passphrase
}

func New(opts ...Option) (*Keystore, error) {
	k := &Keystore{
		id:      uuid.New().String(),
		version: 1,
		entries: []entry.Entry{},
	}

	for _, opt := range opts {
		err := opt(k)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	// TODO: remove this
	k.id = "0000000000000000000000"

	return k, nil
}
