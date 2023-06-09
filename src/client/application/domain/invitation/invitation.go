package invitation

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"myst/pkg/uuid"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/user"
)

var (
	// ErrAccepted    = errors.New("invitation already accepted")
	// ErrNotAccepted = errors.New("invitation not accepted")
	// ErrFinalized   = errors.New("invitation already finalized")
	ErrCannotAccept   = errors.New("cannot accept non-pending invitation")
	ErrCannotFinalize = errors.New("cannot finalize non-accepted invitation")
)

type Invitation struct {
	Id                   string
	Keystore             keystore.Keystore
	Inviter              user.User
	Invitee              user.User
	EncryptedKeystoreKey []byte
	Status               Status
	CreatedAt            time.Time
	UpdatedAt            time.Time
	AcceptedAt           time.Time
	DeclinedAt           time.Time
	DeletedAt            time.Time
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
		i.Keystore, i.Inviter, i.Invitee,
		i.Status,
		base64.StdEncoding.EncodeToString(i.EncryptedKeystoreKey),
		i.CreatedAt, i.UpdatedAt,
		i.AcceptedAt, i.DeclinedAt, i.DeletedAt,
	)
}
