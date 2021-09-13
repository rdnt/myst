package keystore

type Option func(*Keystore) error

func WithName(name string) Option {
	return func(k *Keystore) error {
		k.name = name
		return nil
	}
}

func WithKeystore(keystore []byte) Option {
	return func(k *Keystore) error {
		k.keystore = keystore
		return nil
	}
}
