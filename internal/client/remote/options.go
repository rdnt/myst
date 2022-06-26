package remote

type Option func(*remote)

func WithAddress(address string) Option {
	return func(r *remote) {
		r.address = address
	}
}
