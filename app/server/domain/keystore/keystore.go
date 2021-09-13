package keystore

import (
	"myst/app/server/domain/user"
	"myst/pkg/logger"
	"myst/pkg/timestamp"
	"myst/pkg/uuid"
)

type Keystore struct {
	id        string
	name      string
	keystore  []byte
	owner     *user.User
	createdAt timestamp.Timestamp
	updatedAt timestamp.Timestamp
}

func (k *Keystore) Id() string {
	return k.id
}

func (k *Keystore) Name() string {
	return k.name
}

func (k *Keystore) Owner() string {
	return k.name
}

func (k *Keystore) SetOwner(u *user.User) {
	k.owner = u
}

func (k *Keystore) SetName(name string) {
	k.name = name
	k.updatedAt = timestamp.New()
}

func (k *Keystore) Keystore() []byte {
	return k.keystore
}

func (k *Keystore) SetKeystore(keystore []byte) {
	k.keystore = keystore
	k.updatedAt = timestamp.New()
}

func (k *Keystore) CreatedAt() timestamp.Timestamp {
	return k.createdAt
}

func (k *Keystore) UpdatedAt() timestamp.Timestamp {
	return k.updatedAt
}

func New(opts ...Option) (*Keystore, error) {
	k := &Keystore{
		id:        uuid.New().String(),
		createdAt: timestamp.New(),
		updatedAt: timestamp.New(),
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
