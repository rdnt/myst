package keystoreservice

type Option func(s *service) error

func WithKeystoreRepository(repo KeystoreRepository) Option {
	return func(s *service) error {
		s.keystoreRepo = repo
		return nil
	}
}
