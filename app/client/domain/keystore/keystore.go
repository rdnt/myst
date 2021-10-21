package keystore

import (
	"errors"

	"myst/app/client/domain/keystore/entry"
	"myst/pkg/logger"
	"myst/pkg/uuid"
)

var (
	ErrNotFound       = errors.New("keystore not found")
	ErrEntryNotFound  = errors.New("entry not found")
	ErrEntryExists    = errors.New("entry already exists")
	ErrInvalidEntries = errors.New("invalid entries")
)

type Keystore struct {
	id      string
	version int
	entries []*entry.Entry
}

func (k *Keystore) Id() string {
	return k.id
}

func (k *Keystore) Version() int {
	return k.version
}

func (k *Keystore) SetVersion(version int) {
	k.version = version
}

func (k *Keystore) Entries() []*entry.Entry {
	return k.entries
}

func (k *Keystore) SetEntries(entries []*entry.Entry) error {
	if entries == nil {
		return ErrInvalidEntries
	}

	k.entries = entries
	return nil
}

func (k *Keystore) AddEntry(entry *entry.Entry) error {
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

func New(opts ...Option) (*Keystore, error) {
	k := &Keystore{
		id:      uuid.New().String(),
		version: 1,
		entries: make([]*entry.Entry, 0),
	}

	for _, opt := range opts {
		err := opt(k)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	return k, nil
}
