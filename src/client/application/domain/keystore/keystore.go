package keystore

import (
	"errors"
	"time"

	"myst/pkg/uuid"
	"myst/src/client/application/domain/keystore/entry"
)

var (
	ErrNotFound      = errors.New("keystore not found")
	ErrEntryNotFound = errors.New("entry not found")
	ErrEntryExists   = errors.New("entry already exists")
)

type Keystore struct {
	Id       string
	RemoteId string
	Name     string
	Version  int
	Entries  map[string]entry.Entry
	Access   string

	CreatedAt time.Time
	UpdatedAt time.Time

	Key []byte
}

func (k *Keystore) IncrementVersion() {
	k.Version++
}

func New(opts ...Option) Keystore {
	k := Keystore{
		Id:      uuid.New().String(),
		Version: 1,
		Entries: map[string]entry.Entry{},
		Access:  "read/write",
	}

	for _, opt := range opts {
		opt(&k)
	}

	return k
}
