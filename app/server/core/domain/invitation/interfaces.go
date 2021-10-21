package invitation

import (
	"errors"
)

var (
	ErrNotFound = errors.New("invitation not found")
)

type Repository interface {
	Create(opts ...Option) (*Invitation, error)
	Invitation(id string) (*Invitation, error)
	Update(k *Invitation) error
	Invitations() ([]*Invitation, error)
	Delete(id string) error
}

type Service interface {
	Create(opts ...Option) (*Invitation, error)
}
