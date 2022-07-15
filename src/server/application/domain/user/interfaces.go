package user

import (
	"errors"
)

var (
	ErrNotFound = errors.New("user not found")
)

type Repository interface {
	CreateUser(opts ...Option) (User, error)
	User(id string) (User, error)
	UserByUsername(username string) (User, error)
	UpdateUser(*User) error
	Users() ([]User, error)
	DeleteUser(id string) error
	VerifyPassword(id, password string) (bool, error)
}

type Service interface {
	CreateUser(username, password string, publicKey []byte) (User, error)
	AuthorizeUser(username, password string) error
	User(id string) (User, error)
	UserByUsername(username string) (User, error)
}
