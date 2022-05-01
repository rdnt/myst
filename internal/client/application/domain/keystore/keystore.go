package keystore

import (
	"errors"

	"myst/internal/client/application/domain/entry"

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
	entries  map[string]entry.Entry
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

func (k *Keystore) Entries() map[string]entry.Entry {
	return k.entries
}

func (k *Keystore) SetEntries(entries map[string]entry.Entry) {
	k.entries = entries
}

//func (k *Keystore) CreateEntry(opts ...entry.Option) error {
//	e := entry.New(opts...)
//
//	if _, ok := k.entries[e.Id()]; ok {
//		return ErrEntryExists
//	}
//
//	k.entries[e.Id()] = e
//
//	return nil
//}
//
//func (k *Keystore) DeleteEntry(id string) error {
//	if _, ok := k.entries[id]; !ok {
//		return ErrEntryNotFound
//	}
//
//	delete(k.entries, id)
//
//	return nil
//}

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
		entries: map[string]entry.Entry{},
	}

	for _, opt := range opts {
		opt(k)
	}

	// TODO: remove this
	//k.id = "0000000000000000000000"

	return k
}
