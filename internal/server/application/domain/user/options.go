package user

type Option func(*User)

func WithUsername(username string) Option {
	return func(u *User) {
		u.Username = username
	}
}

func WithPassword(password string) Option {
	return func(u *User) {
		u.Password = password
	}
}

func WithPublicKey(publicKey []byte) Option {
	return func(u *User) {
		u.PublicKey = publicKey
	}
}
