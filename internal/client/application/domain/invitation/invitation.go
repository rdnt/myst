package invitation

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

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
	KeystoreId  string
	InviterId   string
	InviteeId   string
	InviterKey  []byte
	InviteeKey  []byte
	KeystoreKey []byte
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func New(opts ...Option) Invitation {
	inv := Invitation{
		Id:        uuid.New().String(),
		Status:    Pending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		opt(&inv)
	}

	return inv
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
	return fmt.Sprintln(i.Id, i.KeystoreId, i.InviterId, i.InviteeId, i.Status, base64.StdEncoding.EncodeToString(i.InviterKey), base64.StdEncoding.EncodeToString(i.InviteeKey), base64.StdEncoding.EncodeToString(i.KeystoreKey))
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
