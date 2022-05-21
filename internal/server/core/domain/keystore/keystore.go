package keystore

import (
	"fmt"
	"time"

	"myst/pkg/logger"
	"myst/pkg/uuid"
)

type Keystore struct {
	id        string
	name      string
	payload   []byte
	ownerId   string
	viewerIds []string
	createdAt time.Time
	updatedAt time.Time
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
}

func (k *Keystore) Payload() []byte {
	return k.payload
}

func (k *Keystore) SetPayload(payload []byte) {
	k.payload = payload
}

func (k *Keystore) CreatedAt() time.Time {
	return k.createdAt
}

func (k *Keystore) UpdatedAt() time.Time {
	return k.updatedAt
}

func New(opts ...Option) (*Keystore, error) {
	k := &Keystore{
		id:        uuid.New().String(),
		createdAt: time.Now(),
		updatedAt: time.Now(),
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
