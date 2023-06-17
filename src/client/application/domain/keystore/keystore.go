package keystore

import (
	"fmt"
	"time"

	"myst/pkg/uuid"
	"myst/src/client/application/domain/entry"
)

type Keystore struct {
	Id string
	// RemoteId is the id of the keystore on the remote server.
	RemoteId string
	Name     string
	// ReadOnly indicates if the user is allowed to modify the keystore.
	ReadOnly bool
	// Version is bumped every time the keystore is updated, allowing to sync
	// the keystore against the remote server.
	Version int
	// Entries contain the actual website entries of the keystore.
	Entries map[string]entry.Entry
	// Key is the encryption key of a keystore. The keystore itself is always
	// sent encrypted to the remote server and the key is itself end-to-end
	// encrypted before being transferred to others.
	// TODO: what happens if we revoke access to a user? We need to change the
	//  key and update it to all other invitations we have with other users. So
	//  at the end of the day we should be able to remember which users we have
	//  marked as verified to do this silently.
	Key []byte

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (k Keystore) String() string {
	return fmt.Sprintf("Keystore{Id: %s, Name: %s, Version: %d, Entries: %d, Key: %s, CreatedAt: %s, UpdatedAt: %s}",
		k.Id, k.Name, k.Version, len(k.Entries), k.Key, k.CreatedAt, k.UpdatedAt)
}

func New(opts ...Option) Keystore {
	k := Keystore{
		Id:      uuid.New().String(),
		Version: 1,
		Entries: map[string]entry.Entry{},
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&k)
		}
	}

	return k
}

type Option func(*Keystore)

func WithName(name string) Option {
	return func(k *Keystore) {
		k.Name = name
	}
}
