package keystore

type Option func(*Keystore)

func WithName(name string) Option {
	return func(k *Keystore) {
		k.Name = name
	}
}

func WithPayload(payload []byte) Option {
	return func(k *Keystore) {
		k.Payload = payload
	}
}

func WithOwnerId(id string) Option {
	return func(k *Keystore) {
		k.OwnerId = id
	}
}
