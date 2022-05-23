package entry

type Option func(*Entry)

func WithId(id string) Option {
	return func(e *Entry) {
		e.Id = id
	}
}

func WithWebsite(website string) Option {
	return func(e *Entry) {
		e.Website = website
	}
}

func WithUsername(username string) Option {
	return func(e *Entry) {
		e.Username = username
	}
}

func WithPassword(password string) Option {
	return func(e *Entry) {
		e.Password = password
	}
}

func WithNotes(notes string) Option {
	return func(e *Entry) {
		e.Notes = notes
	}
}
