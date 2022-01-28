package invitation

type Option func(i *Invitation) error

func WithId(id string) Option {
	return func(i *Invitation) error {
		i.id = id
		return nil
	}
}

//func WithInviter(u *user.User) Option {
//	return func(i *Invitation) error {
//		i.inviter = u
//		return nil
//	}
//}
//
//func WithKeystore(k *keystore.Keystore) Option {
//	return func(i *Invitation) error {
//		i.keystore = k
//		return nil
//	}
//}
//
//func WithInvitee(u *user.User) Option {
//	return func(i *Invitation) error {
//		i.invitee = u
//		return nil
//	}
//}
//
//func WithInviterKey(key []byte) Option {
//	return func(i *Invitation) error {
//		i.inviterKey = key
//		return nil
//	}
//}
