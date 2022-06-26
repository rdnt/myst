package remote

type Option func(*remote)

func WithAddress(address string) Option {
	return func(r *remote) {
		r.address = address
	}
}

func WithPublicKey(publicKey []byte) Option {
	return func(r *remote) {
		r.publicKey = publicKey
	}
}

func WithPrivateKey(privateKey []byte) Option {
	return func(r *remote) {
		r.privateKey = privateKey
	}
}
