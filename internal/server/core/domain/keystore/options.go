package keystore

type Option func(*Keystore) error

func WithName(name string) Option {
	return func(k *Keystore) error {
		//k.Id = name // TODO: remove
		k.Name = name
		return nil
	}
}

func WithPayload(payload []byte) Option {
	return func(k *Keystore) error {
		k.Payload = payload
		return nil
	}
}

func WithOwnerId(id string) Option {
	return func(k *Keystore) error {
		k.OwnerId = id
		return nil
	}
}
