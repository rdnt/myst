package invitation

type Option func(i *Invitation)

func WithInviterId(id string) Option {
	return func(i *Invitation) {
		i.InviterId = id
	}
}

func WithKeystoreId(id string) Option {
	return func(i *Invitation) {
		i.KeystoreId = id
	}
}

func WithInviteeId(id string) Option {
	return func(i *Invitation) {
		i.InviteeId = id
	}
}

func WithInviterKey(key []byte) Option {
	return func(i *Invitation) {
		i.InviterKey = key
	}
}
