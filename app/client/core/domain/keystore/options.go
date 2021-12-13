package keystore

type Option func(*Keystore) error

func WithName(name string) Option {
	return func(k *Keystore) error {
		k.name = name
		return nil
	}
}

func WithPassphrase(passphrase []byte) Option {
	return func(k *Keystore) error {
		k.passphrase = passphrase
		return nil
	}
}
