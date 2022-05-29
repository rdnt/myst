package remote

type Option func(*remote)

func WithAddress(address string) Option {
	return func(r *remote) {
		r.address = address
	}
}

func WithUsername(username string) Option {
	return func(r *remote) {
		r.username = username
	}
}

func WithPassword(password string) Option {
	return func(r *remote) {
		r.password = password
	}
}
