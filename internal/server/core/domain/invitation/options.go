package invitation

type Option func(i *Invitation) error

func WithInviterId(id string) Option {
	return func(i *Invitation) error {
		i.inviterId = id
		return nil
	}
}

func WithKeystoreId(id string) Option {
	return func(i *Invitation) error {
		i.keystoreId = id
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
