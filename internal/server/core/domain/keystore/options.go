package keystore

type Option func(*Keystore) error

func WithName(name string) Option {
	return func(k *Keystore) error {
		k.id = name // TODO: remove
		k.name = name
		return nil
	}
}

func WithPayload(payload []byte) Option {
	return func(k *Keystore) error {
		k.payload = payload
		return nil
	}
}

func WithOwnerId(id string) Option {
	return func(k *Keystore) error {
		k.ownerId = id
		return nil
	}
}
