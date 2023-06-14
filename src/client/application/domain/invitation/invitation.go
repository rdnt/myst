package invitation

import (
	"encoding/base64"
	"fmt"
	"time"

	"myst/src/client/application/domain/user"
)

type Invitation struct {
	Id                   string
	RemoteKeystoreId     string
	KeystoreName         string
	Inviter              user.User
	Invitee              user.User
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
		i.RemoteKeystoreId, i.KeystoreName,
		i.Inviter.Id, i.Invitee.Id,
		i.Status,
		base64.StdEncoding.EncodeToString(i.EncryptedKeystoreKey),
		i.CreatedAt, i.UpdatedAt,
		i.AcceptedAt, i.DeclinedAt, i.DeletedAt,
	)
}
