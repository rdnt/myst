package invitation

import (
	"errors"
)

type Status string

const (
	Pending   Status = "pending"
	Accepted  Status = "accepted"
	Finalized Status = "finalized"
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
	case Finalized:
		return Finalized, nil
	default:
		return "", errors.New("invalid status")
	}
}
