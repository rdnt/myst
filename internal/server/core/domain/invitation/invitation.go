package invitation

import (
	"errors"
	"fmt"

	"myst/pkg/logger"
	"myst/pkg/timestamp"
	"myst/pkg/uuid"
)

var (
	ErrAlreadyAccepted  = errors.New("invitation already accepted")
	ErrNotAccepted      = errors.New("invitation not accepted")
	ErrAlreadyFinalized = errors.New("invitation already finalized")
	ErrCannotAccept     = errors.New("cannot accept non-pending invitation")
	ErrCannotFinalize   = errors.New("cannot finalize non-accepted invitation")
)

type Invitation struct {
	id          string
	inviterId   string
	keystoreId  string
	inviteeId   string
	inviterKey  []byte
	inviteeKey  []byte
	keystoreKey []byte // encrypted
	status      Status
	createdAt   timestamp.Timestamp
	updatedAt   timestamp.Timestamp
}

func New(opts ...Option) (*Invitation, error) {
	k := &Invitation{
		id:        uuid.New().String(),
		createdAt: timestamp.New(),
		updatedAt: timestamp.New(),
		status:    Pending,
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

func (i *Invitation) Id() string {
	return i.id
}

func (i *Invitation) InviterId() string {
	return i.inviterId
}

func (i *Invitation) KeystoreId() string {
	return i.keystoreId
}

func (i *Invitation) InviteeId() string {
	return i.inviteeId
}

func (i *Invitation) InviterKey() []byte {
	return i.inviterKey
}

func (i *Invitation) InviteeKey() []byte {
	return i.inviteeKey
}

func (i *Invitation) KeystoreKey() []byte {
	return i.keystoreKey
}

func (i *Invitation) Pending() bool {
	return i.status == Pending
}

func (i *Invitation) Accepted() bool {
	return i.status == Accepted
}

func (i *Invitation) Finalized() bool {
	return i.status == Finalized
}

func (i *Invitation) CreatedAt() timestamp.Timestamp {
	return i.createdAt
}

func (i *Invitation) UpdatedAt() timestamp.Timestamp {
	return i.updatedAt
}

func (i *Invitation) Status() Status {
	return i.status
}

func (i *Invitation) String() string {
	return fmt.Sprintln(i.id, i.inviterId, i.keystoreId, i.inviteeId, i.status)
}

func (i *Invitation) Accept(inviteeKey []byte) error {
	if i.status != Pending {
		return ErrCannotAccept
	}

	i.inviteeKey = inviteeKey
	i.status = Accepted

	return nil
}

func (i *Invitation) Finalize(keystoreKey []byte) error {
	if i.status != Accepted {
		return ErrCannotFinalize
	}

	i.keystoreKey = keystoreKey
	i.status = Finalized

	return nil
}
