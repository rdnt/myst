package keystore

import (
	"fmt"
	"time"

	"myst/pkg/logger"
	"myst/pkg/uuid"
)

type Keystore struct {
	Id        string
	Name      string
	Payload   []byte
	OwnerId   string
	ViewerIds []string // TODO: can we make do without storing users with access on the keystore itself?
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (k *Keystore) String() string {
	return fmt.Sprintln(k.Id, k.Name, k.OwnerId, k.ViewerIds, k.CreatedAt, k.UpdatedAt)
}

func (k *Keystore) SetOwnerId(id string) {
	k.OwnerId = id
}

func (k *Keystore) SetName(name string) {
	k.Name = name
}

func (k *Keystore) SetPayload(payload []byte) {
	k.Payload = payload
}

func New(opts ...Option) (Keystore, error) {
	k := Keystore{
		Id:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		err := opt(&k)
		if err != nil {
			logger.Error(err)
			return Keystore{}, err
		}
	}

	// TODO: remove this
	//k.Id = "0000000000000000000000"

	return k, nil
}
