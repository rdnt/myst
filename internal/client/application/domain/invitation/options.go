package invitation

type Option func(i *Invitation) error

func WithInviterId(id string) Option {
	return func(i *Invitation) error {
		i.InviterId = id
		return nil
	}
}

func WithKeystoreId(id string) Option {
	return func(i *Invitation) error {
		i.KeystoreId = id
		return nil
	}
}

func WithInviteeId(id string) Option {
	return func(i *Invitation) error {
		i.InviteeId = id
		return nil
	}
}

func WithInviterKey(key []byte) Option {
	return func(i *Invitation) error {
		i.InviterKey = key
		return nil
	}
}
