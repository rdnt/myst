package user

import (
	"errors"
	"fmt"
	"time"

	"myst/pkg/uuid"
)

var (
	ErrInvalidUsername = errors.New("invalid username")
)

type User struct {
	Id        string
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) String() string {
	return fmt.Sprintln(u.Id, u.Username, "****", u.CreatedAt, u.UpdatedAt)
}

func New(opts ...Option) (User, error) {
	u := User{
		Id:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		opt(&u)
	}

	if u.Username == "" {
		return User{}, ErrInvalidUsername
	}

	return u, nil
}
