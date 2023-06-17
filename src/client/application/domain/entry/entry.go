package entry

import (
	"time"

	"myst/pkg/uuid"
)

type Entry struct {
	Id        string
	Website   string
	Username  string
	Password  string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(opts ...Option) Entry {
	e := Entry{
		Id:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&e)
		}
	}

	return e
}

type Option func(*Entry)

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
