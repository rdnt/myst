package invitation

import (
	"myst/internal/server/core/domain/keystore"
	"myst/internal/server/core/domain/user"
)

type Option func(i *Invitation) error

func WithInviter(u user.User) Option {
	return func(i *Invitation) error {
		i.inviter = u
		return nil
	}
}

func WithKeystore(k keystore.Keystore) Option {
	return func(i *Invitation) error {
		i.keystore = k
		return nil
	}
}

func WithInvitee(u user.User) Option {
	return func(i *Invitation) error {
		i.invitee = u
		return nil
	}
}

func WithInviterKey(key []byte) Option {
	return func(i *Invitation) error {
		i.inviterKey = key
		return nil
	}
}
