package repository

import (
	"encoding/json"
	"os"
	"sync"

	"myst/src/server/application/domain/invitation"
	"myst/src/server/application/domain/keystore"
)

type Repository struct {
	mux         sync.Mutex
	invitations map[string]invitation.Invitation
	keystores   map[string]keystore.Keystore
	users       map[string]User
}

func New() *Repository {
	r := &Repository{
		keystores:   map[string]keystore.Keystore{},
		invitations: map[string]invitation.Invitation{},
		users:       map[string]User{},
	}

	{ // debug
		b, err := os.ReadFile("keystores.json")
		if err == nil {
			_ = json.Unmarshal(b, &r.keystores)
		}

		b, err = os.ReadFile("users.json")
		if err == nil {
			_ = json.Unmarshal(b, &r.users)
		}

		b, err = os.ReadFile("invitations.json")
		if err == nil {
			_ = json.Unmarshal(b, &r.invitations)
		}
	}

	return r
}
