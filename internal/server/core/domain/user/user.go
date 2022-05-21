package user

import (
	"errors"
	"fmt"
	"time"

	"myst/pkg/logger"
	"myst/pkg/uuid"
)

var (
	ErrInvalidUsername = errors.New("invalid username")
)

type User struct {
	id        string
	username  string
	password  string
	createdAt time.Time
	updatedAt time.Time
}

func (u *User) Id() string {
	return u.id
}

func (u *User) String() string {
	return fmt.Sprintln(u.id, u.username, u.password)
}

func (u *User) Username() string {
	return u.username
}

func (u *User) Password() string {
	return u.password
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func New(opts ...Option) (*User, error) {
	u := &User{
		id:        uuid.New().String(),
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	for _, opt := range opts {
		err := opt(u)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	if u.username == "" {
		return nil, ErrInvalidUsername
	}

	return u, nil
}
