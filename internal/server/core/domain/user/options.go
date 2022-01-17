package user

type Option func(*User) error

func WithUsername(username string) Option {
	return func(u *User) error {
		u.id = username // TODO: remove
		u.username = username
		return nil
	}
}

func WithPassword(password string) Option {
	return func(u *User) error {
		u.password = password
		return nil
	}
}
