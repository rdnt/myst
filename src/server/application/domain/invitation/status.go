package invitation

import (
	"errors"
)

type Status string

const (
	Pending   Status = "pending"
	Deleted   Status = "deleted"
	Declined  Status = "declined"
	Accepted  Status = "accepted"
	Cancelled Status = "cancelled"
	Finalized Status = "finalized"
)

func (s Status) String() string {
	return string(s)
}

func StatusFromString(s string) (Status, error) {
	switch Status(s) {
	case Pending:
		return Pending, nil
	case Deleted:
		return Deleted, nil
	case Declined:
		return Declined, nil
	case Accepted:
		return Accepted, nil
	case Cancelled:
		return Cancelled, nil
	case Finalized:
		return Finalized, nil
	default:
		return "", errors.New("invalid status")
	}
}
