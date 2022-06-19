package invitation

import (
	"errors"
	"fmt"
	"time"

	"myst/pkg/logger"
	"myst/pkg/uuid"
)

var (
	//ErrAlreadyAccepted  = errors.New("invitation already accepted")
	//ErrNotAccepted      = errors.New("invitation not accepted")
	//ErrAlreadyFinalized = errors.New("invitation already finalized")
	ErrCannotAccept   = errors.New("cannot accept non-pending invitation")
	ErrCannotFinalize = errors.New("cannot finalize non-accepted invitation")
	ErrCannotDecline  = errors.New("cannot decline non-pending invitation")
	ErrCannotCancel   = errors.New("cannot cancel non-pending invitation")
)

type Invitation struct {
	Id           string
	InviterId    string
	KeystoreId   string
	KeystoreName string
	InviteeId    string
	InviterKey   []byte
	InviteeKey   []byte
	KeystoreKey  []byte
	Status       Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
	AcceptedAt   time.Time
	DeclinedAt   time.Time
	DeletedAt    time.Time
}

func New(opts ...Option) (*Invitation, error) {
	k := &Invitation{
		Id:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    Pending,
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

func (i *Invitation) String() string {
	return fmt.Sprintln(i.Id, i.InviteeKey, i.KeystoreKey, i.InviteeId, i.Status)
}

func (i *Invitation) Accept(inviteeKey []byte) error {
	if i.Status != Pending {
		return ErrCannotAccept
	}

	i.InviteeKey = inviteeKey
	i.Status = Accepted
	i.AcceptedAt = time.Now()

	return nil
}

func (i *Invitation) Finalize(keystoreKey []byte) error {
	if i.Status != Accepted {
		return ErrCannotFinalize
	}

	i.KeystoreKey = keystoreKey
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
