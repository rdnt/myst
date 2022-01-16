package user

import (
	"errors"
	"fmt"

	"myst/pkg/logger"
	"myst/pkg/timestamp"
	"myst/pkg/uuid"
)

var (
	ErrInvalidUsername = errors.New("invalid username")
)

type User struct {
	id          string
	username    string
	password    string
	keystoreIds []string
	createdAt   timestamp.Timestamp
	updatedAt   timestamp.Timestamp
}

func (u *User) KeystoreIds() []string {
	return u.keystoreIds
}

func (u *User) SetKeystoreIds(ids []string) {
	u.keystoreIds = ids
}

func (u *User) OwnKeystore(id string) {
	u.keystoreIds = append(u.keystoreIds, id)
}

func (u *User) Id() string {
	return u.id
}

func (u *User) String() string {
	return fmt.Sprintln(u.id, u.username, u.password, u.keystoreIds)
}

func (u *User) Username() string {
	return u.username
}

func (u *User) Password() string {
	return u.password
}

func (u *User) CreatedAt() timestamp.Timestamp {
	return u.createdAt
}

func (u *User) UpdatedAt() timestamp.Timestamp {
	return u.updatedAt
}

func New(opts ...Option) (*User, error) {
	u := &User{
		id:        uuid.New().String(),
		createdAt: timestamp.New(),
		updatedAt: timestamp.New(),
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
