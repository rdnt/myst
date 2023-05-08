package user

type Option func(*User)

func WithUsername(username string) Option {
	return func(u *User) {
		u.Username = username
	}
}

func WithPasswordHash(hash string) Option {
	return func(u *User) {
		u.PasswordHash = hash
	}
}

func WithPublicKey(publicKey []byte) Option {
	return func(u *User) {
		u.PublicKey = publicKey
	}
}
