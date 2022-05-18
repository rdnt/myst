package user

import (
	"errors"
)

var (
	ErrNotFound = errors.New("user not found")
)

type Repository interface {
	CreateUser(opts ...Option) (*User, error)
	User(id string) (*User, error)
	UpdateUser(*User) error
	Users() ([]*User, error)
	DeleteUser(id string) error
}

type Service interface {
	CreateUser(opts ...Option) (*User, error)
	AuthorizeUser(u *User, password string) error
	User(id string) (*User, error)
}
