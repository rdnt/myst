package keystore

import (
	"fmt"
	"time"

	"myst/pkg/uuid"
)

type Keystore struct {
	Id        string
	Name      string
	Payload   []byte
	OwnerId   string
	ViewerIds []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (k Keystore) String() string {
	return fmt.Sprintln(k.Id, k.Name, k.OwnerId, k.ViewerIds, k.CreatedAt, k.UpdatedAt)
}

func New(opts ...Option) Keystore {
	k := Keystore{
		Id:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

func WithPayload(payload []byte) Option {
	return func(k *Keystore) {
		k.Payload = payload
	}
}

func WithOwnerId(id string) Option {
	return func(k *Keystore) {
		k.OwnerId = id
	}
}
