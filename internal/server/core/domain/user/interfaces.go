package user

import "errors"

var (
	ErrNotFound = errors.New("user not found")
)

type Repository interface {
	Create(opts ...Option) (*User, error)
	User(id string) (*User, error)
	Update(*User) error
	Users() ([]*User, error)
	Delete(id string) error
}

type Service interface {
	Register(opts ...Option) (*User, error)
	Authorize(u *User, password string) error
}
