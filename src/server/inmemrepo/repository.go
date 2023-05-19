package inmemrepo

import (
	"sync"

	"myst/src/server/application/domain/invitation"
	"myst/src/server/application/domain/keystore"
	"myst/src/server/application/domain/user"
)

type Repository struct {
	mux         sync.Mutex
	invitations map[string]invitation.Invitation
	keystores   map[string]keystore.Keystore
	users       map[string]user.User
}

func New() *Repository {
	r := &Repository{
		keystores:   map[string]keystore.Keystore{},
		invitations: map[string]invitation.Invitation{},
		users:       map[string]user.User{},
	}

	return r
}
