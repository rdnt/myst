package invitation

import (
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/user"
)

type Option func(i *Invitation)

func WithKeystore(k keystore.Keystore) Option {
	return func(i *Invitation) {
		i.Keystore = k
	}
}

func WithInviteeUsername(username string) Option {
	return func(i *Invitation) {
		i.Invitee = user.User{Username: username}
	}
}
