package user

import (
	"encoding/base64"
	"fmt"
)

type User struct {
	Id        string
	Username  string
	PublicKey []byte
	// TODO: do not expose the shared secret directly, hash it and store
	// the hash internally on the enclave for each user that has an invitation
	// accepted/finalized.
	SharedSecret []byte
	Verified     bool
}

func (u User) String() string {
	return fmt.Sprintln(
		u.Id, u.Username,
		base64.StdEncoding.EncodeToString(u.PublicKey),
		u.Verified,
	)
}
