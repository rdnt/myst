package invitation

import (
	"errors"
)

type Status string

const (
	Pending         Status = "pending"
	Accepted        Status = "accepted"
	Declined        Status = "declined"
	Deleted         Status = "deleted"
	InviterVerified Status = "inviter-verified"
	InviteeVerified Status = "invitee-verified"
	Finalized       Status = "finalized"
)

func (s Status) String() string {
	return string(s)
}

func StatusFromString(s string) (Status, error) {
	switch Status(s) {
	case Pending:
		return Pending, nil
	case Accepted:
		return Accepted, nil
	case Declined:
		return Declined, nil
	case Deleted:
		return Deleted, nil
	case InviterVerified:
		return InviterVerified, nil
	case InviteeVerified:
		return InviteeVerified, nil
	case Finalized:
		return Finalized, nil
	default:
		return "", errors.New("invalid status")
	}
}
