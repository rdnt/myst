package keystore

import (
	"errors"
	"time"

	"myst/internal/client/application/domain/entry"
	"myst/pkg/uuid"
)

var (
	ErrNotFound      = errors.New("keystore not found")
	ErrEntryNotFound = errors.New("entry not found")
	ErrEntryExists   = errors.New("entry already exists")
)

type Keystore struct {
	Id      string
	Name    string
	Version int
	Entries map[string]entry.Entry

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (k *Keystore) SetName(name string) {
	k.Name = name
}

func (k *Keystore) IncrementVersion() {
	k.Version++
}

func (k *Keystore) SetEntries(entries map[string]entry.Entry) {
	k.Entries = entries
}

func New(opts ...Option) Keystore {
	k := Keystore{
		Id:      uuid.New().String(),
		Version: 1,
		Entries: map[string]entry.Entry{},
	}

	for _, opt := range opts {
		opt(&k)
	}

	// TODO: remove this
	//k.id = "0000000000000000000000"

	return k
}
