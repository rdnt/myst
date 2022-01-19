package invitation

import (
	"myst/internal/server/core/domain/keystore"
)

type Option func(i *Invitation) error

func WithInviterId(id string) Option {
	return func(i *Invitation) error {
		i.inviterId = id
		return nil
	}
}

func WithKeystore(k *keystore.Keystore) Option {
	return func(i *Invitation) error {
		i.keystore = k
		return nil
	}
}

func WithInviteeId(id string) Option {
	return func(i *Invitation) error {
		i.inviteeId = id
		return nil
	}
}

func WithInviterKey(key []byte) Option {
	return func(i *Invitation) error {
		i.inviterKey = key
		return nil
	}
}
