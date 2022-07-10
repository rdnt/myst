package user

import (
	"encoding/base64"
	"errors"
	"fmt"
)

var (
	ErrInvalidUsername = errors.New("invalid username")
)

type User struct {
	Id        string
	Username  string
	PublicKey []byte
}

func (u User) String() string {
	return fmt.Sprintln(
		u.Id, u.Username,
		base64.StdEncoding.EncodeToString(u.PublicKey),
	)
}
