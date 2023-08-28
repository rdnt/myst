package invitation

import (
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
	k := Invitation{
		Id:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    Pending,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&k)
		}
	}

	return k
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

func (i Invitation) String() string {
	return fmt.Sprintf(
		"Invitation{Id: %s, KeystoreId: %s, InviterId: %s, InviteeId: %s, EncryptedKeystoreKey: %s, Status: %s}",
		i.Id,
		i.KeystoreId, i.InviterId, i.InviteeId,
		i.EncryptedKeystoreKey,
		i.Status,
	)
}

func (i Invitation) Pending() bool {
	return i.Status == Pending
}

func (i Invitation) Accepted() bool {
	return i.Status == Accepted
}

func (i Invitation) Deleted() bool {
	return i.Status == Deleted
}

func (i Invitation) Declined() bool {
	return i.Status == Declined
}

func (i Invitation) Finalized() bool {
	return i.Status == Finalized
}
