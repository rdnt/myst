package invitation

type Option func(i *Invitation)

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
