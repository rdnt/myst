package user

type Option func(*User)

func WithUsername(username string) Option {
	return func(u *User) {
		u.Id = username // TODO: remove
		u.Username = username
	}
}

func WithPassword(password string) Option {
	return func(u *User) {
		u.Password = password
	}
}
