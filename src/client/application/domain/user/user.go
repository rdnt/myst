package user

import (
	"encoding/base64"
	"fmt"
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
