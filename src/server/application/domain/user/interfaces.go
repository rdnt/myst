package user

import (
	"errors"
)

var (
	ErrNotFound = errors.New("user not found")
)

// type Repository interface {
// 	CreateUser(User) (User, error)
// 	User(id string) (User, error)
// 	UserByUsername(username string) (User, error)
// 	UpdateUser(User) (User, error)
// 	Users() ([]User, error)
// 	// DeleteUser(id string) error
// }

// type Service interface {
// 	CreateUser(username, password string, publicKey []byte) (User, error)
// 	AuthorizeUser(username, password string) (User, error)
// 	User(id string) (User, error)
// 	UserByUsername(username string) (User, error)
// }
