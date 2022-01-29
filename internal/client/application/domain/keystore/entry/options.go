package entry

type Option func(*Entry)

func WithId(id string) Option {
	return func(e *Entry) {
		e.id = id
	}
}

func WithLabel(label string) Option {
	return func(e *Entry) {
		e.label = label
	}
}

func WithUsername(username string) Option {
	return func(e *Entry) {
		e.username = username
	}
}

func WithPassword(password string) Option {
	return func(e *Entry) {
		e.password = password
	}
}
