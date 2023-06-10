package invitation

import (
	"encoding/base64"
	"fmt"
	"time"

	"myst/pkg/uuid"
)

type Invitation struct {
	Id                   string
	KeystoreId           string
	InviterId            string
	InviteeId            string
	EncryptedKeystoreKey []byte
	Status               Status
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            time.Time
	DeclinedAt           time.Time
	AcceptedAt           time.Time
	CancelledAt          time.Time
	FinalizedAt          time.Time
}

func New(opts ...Option) Invitation {
	inv := Invitation{
		Id:        uuid.New().String(),
		Status:    Pending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&inv)
		}
	}

	return inv
}

type Option func(i *Invitation)

func WithKeystoreId(id string) Option {
	return func(i *Invitation) {
		i.KeystoreId = id
	}
}

func WithInviterId(id string) Option {
	return func(i *Invitation) {
		i.InviterId = id
	}
}

func WithInviteeId(id string) Option {
	return func(i *Invitation) {
		i.InviteeId = id
	}
}

func (i Invitation) Pending() bool {
	return i.Status == Pending
}

func (i Invitation) Accepted() bool {
	return i.Status == Accepted
}

func (i Invitation) Finalized() bool {
	return i.Status == Finalized
}

func (i Invitation) String() string {
	return fmt.Sprintln(
		i.Id,
		i.KeystoreId, i.InviterId, i.InviteeId,
		i.Status,
		base64.StdEncoding.EncodeToString(i.EncryptedKeystoreKey),
		i.CreatedAt, i.UpdatedAt,
		i.AcceptedAt, i.DeclinedAt, i.DeletedAt,
	)
}
