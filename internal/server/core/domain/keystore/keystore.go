package keystore

import (
	"fmt"

	"myst/pkg/logger"
	"myst/pkg/timestamp"
	"myst/pkg/uuid"
)

type Keystore struct {
	id        string
	name      string
	payload   []byte
	ownerId   string
	viewerIds []string
	createdAt timestamp.Timestamp
	updatedAt timestamp.Timestamp
}

func (k *Keystore) Id() string {
	return k.id
}

func (k *Keystore) String() string {
	return fmt.Sprintln(k.id, k.name, k.ownerId)
}

func (k *Keystore) ViewerIds() []string {
	return k.viewerIds
}

func (k *Keystore) Name() string {
	return k.name
}

func (k *Keystore) OwnerId() string {
	return k.ownerId
}

func (k *Keystore) SetOwnerId(id string) {
	k.ownerId = id
}

func (k *Keystore) SetName(name string) {
	k.name = name
	k.updatedAt = timestamp.New()
}

func (k *Keystore) Payload() []byte {
	return k.payload
}

func (k *Keystore) SetPayload(payload []byte) {
	k.payload = payload
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

	// TODO: remove this
	k.id = "0000000000000000000000"

	return k, nil
}
