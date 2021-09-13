package user

type Option func(*User) error

func WithUsername(username string) Option {
	return func(u *User) error {
		u.username = username
		return nil
	}
}
