package memdb

import (
	"sync"

	"myst/internal/server/application/domain/invitation"
	"myst/internal/server/application/domain/keystore"
)

type Repository struct {
	mux         sync.Mutex
	invitations map[string]invitation.Invitation
	keystores   map[string]keystore.Keystore
	users       map[string]User
}

func New() *Repository {
	return &Repository{
		keystores:   map[string]keystore.Keystore{},
		invitations: map[string]invitation.Invitation{},
		users:       map[string]User{},
	}
}
