package invitation

import (
	"errors"
	"fmt"
	"time"

	"myst/pkg/uuid"
)

var (
	ErrCannotAccept   = errors.New("cannot accept non-pending invitation")
	ErrCannotFinalize = errors.New("cannot finalize non-accepted invitation")
	ErrCannotDecline  = errors.New("cannot decline non-pending invitation")
	ErrCannotCancel   = errors.New("cannot cancel non-pending invitation")
)

type Invitation struct {
	Id                   string
	KeystoreId           string
	InviterId            string
	InviteeId            string
	EncryptedKeystoreKey []byte
	Status               Status
	InviterVerified      bool
	InviteeVerified      bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
	AcceptedAt           time.Time
	DeclinedAt           time.Time
	DeletedAt            time.Time
}

func New(opts ...Option) Invitation {
	k := Invitation{
		Id:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    Pending,
	}

	for _, opt := range opts {
		opt(&k)
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

func (i *Invitation) String() string {
	return fmt.Sprintln(
		i.Id,
		i.KeystoreId, i.InviterId, i.InviteeId,
		i.EncryptedKeystoreKey,
		i.Status,
	)
}

func (i *Invitation) Pending() bool {
	return i.Status == Pending
}

func (i *Invitation) Accepted() bool {
	return i.Status == Accepted
}

func (i *Invitation) Finalized() bool {
	return i.Status == Finalized
}

func (i *Invitation) Declined() bool {
	return i.Status == Declined
}

func (i *Invitation) Deleted() bool {
	return i.Status == Deleted
}

func (i *Invitation) Accept() error {
	if i.Status != Pending {
		return ErrCannotAccept
	}

	i.Status = Accepted
	i.AcceptedAt = time.Now()

	return nil
}

func (i *Invitation) Finalize(encryptedKeystoreKey []byte) error {
	if i.Status != Accepted {
		return ErrCannotFinalize
	}

	i.EncryptedKeystoreKey = encryptedKeystoreKey
	i.Status = Finalized

	return nil
}

func (i *Invitation) Decline() error {
	if i.Status != Pending {
		return ErrCannotDecline
	}

	i.Status = Declined
	i.DeclinedAt = time.Now()

	return nil
}

func (i *Invitation) Delete() error {
	if i.Status != Pending {
		return ErrCannotCancel
	}

	i.Status = Deleted
	i.DeletedAt = time.Now()

	return nil
}
