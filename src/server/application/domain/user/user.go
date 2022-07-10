package user

import (
	"encoding/hex"
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
	PublicKey []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) String() string {
	return fmt.Sprintln(u.Id, u.Username, "****", hex.EncodeToString(u.PublicKey))
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