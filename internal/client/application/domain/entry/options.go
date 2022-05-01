package entry

type Option func(*Entry)

func WithId(id string) Option {
	return func(e *Entry) {
		e.id = id
	}
}

func WithWebsite(website string) Option {
	return func(e *Entry) {
		e.website = website
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

func WithNotes(notes string) Option {
	return func(e *Entry) {
		e.notes = notes
	}
}
