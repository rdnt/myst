package user

import (
	"encoding/hex"
	"fmt"
	"time"

	"myst/pkg/uuid"
)

type User struct {
	Id           string
	Username     string
	PasswordHash string
	PublicKey    []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u User) String() string {
	return fmt.Sprintln(u.Id, u.Username, "****", hex.EncodeToString(u.PublicKey))
}

func New(opts ...Option) User {
	u := User{
		Id:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&u)
		}
	}

	return u
}
