package user

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidUsername = errors.New("invalid username")
)

type User struct {
	Id       string
	Username string
}

func (u *User) String() string {
	return fmt.Sprintln(u.Id, u.Username)
}
