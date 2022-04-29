package keystore

import (
	"errors"

	"myst/internal/client/application/domain/keystore/entry"

	"myst/pkg/uuid"
)

var (
	ErrNotFound      = errors.New("keystore not found")
	ErrEntryNotFound = errors.New("entry not found")
	ErrEntryExists   = errors.New("entry already exists")
)

type Keystore struct {
	id       string
	name     string
	version  int
	entries  []entry.Entry
	password string
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

type UpdateEntryOptions struct {
	Password *string
	Notes    *string
}

func (k *Keystore) UpdateEntry(id string, opts UpdateEntryOptions) error {
	var entry *entry.Entry
	var index = -1

	for i, e := range k.entries {
		if e.Id() == id {
			entry = &e
			index = i
			break
		}
	}

	if entry == nil {
		return ErrEntryNotFound
	}

	if opts.Password != nil {
		entry.SetPassword(*opts.Password)
	}

	if opts.Notes != nil {
		entry.SetNotes(*opts.Notes)
	}

	k.entries[index] = *entry

	return nil

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

func (k *Keystore) Password() string {
	return k.password
}

func (k *Keystore) SetPassword(password string) {
	k.password = password
}

func New(opts ...Option) *Keystore {
	k := &Keystore{
		id:      uuid.New().String(),
		version: 1,
		entries: []entry.Entry{},
	}

	for _, opt := range opts {
		opt(k)
	}

	// TODO: remove this
	//k.id = "0000000000000000000000"

	return k
}
