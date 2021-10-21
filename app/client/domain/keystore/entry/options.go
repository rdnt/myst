package entry

type Option func(*Entry) error

func WithLabel(label string) Option {
	return func(e *Entry) error {
		e.label = label
		return nil
	}
}

func WithUsername(username string) Option {
	return func(e *Entry) error {
		e.username = username
		return nil
	}
}

func WithPassword(password string) Option {
	return func(e *Entry) error {
		e.password = password
		return nil
	}
}
