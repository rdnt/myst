package invitation

import (
	"errors"
	"fmt"

	"myst/pkg/logger"
	"myst/pkg/timestamp"
	"myst/pkg/uuid"
)

var (
	//ErrAccepted    = errors.New("invitation already accepted")
	//ErrNotAccepted = errors.New("invitation not accepted")
	//ErrFinalized   = errors.New("invitation already finalized")
	ErrCannotAccept   = errors.New("cannot accept non-pending invitation")
	ErrCannotFinalize = errors.New("cannot finalize non-accepted invitation")
)

type Invitation struct {
	Id          string
	InviterId   string
	KeystoreId  string
	InviteeId   string
	InviterKey  []byte
	InviteeKey  []byte
	KeystoreKey []byte
	Status      Status
	CreatedAt   timestamp.Timestamp
	UpdatedAt   timestamp.Timestamp
}

func New(opts ...Option) (*Invitation, error) {
	k := &Invitation{
		Id:        uuid.New().String(),
		CreatedAt: timestamp.New(),
		UpdatedAt: timestamp.New(),
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

func (i *Invitation) String() string {
	return fmt.Sprintln(i.Id, i.InviteeKey, i.KeystoreKey, i.InviteeId, i.Status)
}

func (i *Invitation) Accept(inviteeKey []byte) error {
	if i.Status != Pending {
		return ErrCannotAccept
	}

	i.InviteeKey = inviteeKey
	i.Status = Accepted

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
