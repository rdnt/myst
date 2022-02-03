package entry

import (
	"myst/pkg/uuid"
)

type Entry struct {
	id       string
	website  string
	username string
	password string
	notes    string
}

func (e *Entry) Id() string {
	return e.id
}

func (e *Entry) Website() string {
	return e.website
}

func (e *Entry) SetWebsite(website string) {
	e.website = website
}

func (e *Entry) Username() string {
	return e.username
}

func (e *Entry) SetUsername(username string) {
	e.username = username
}

func (e *Entry) Password() string {
	return e.password
}

func (e *Entry) SetPassword(password string) {
	e.password = password
}

func (e *Entry) Notes() string {
	return e.notes
}

func (e *Entry) SetNotes(notes string) {
	e.notes = notes
}

func New(opts ...Option) Entry {
	e := Entry{
		id: uuid.New().String(),
	}

	for _, opt := range opts {
		opt(&e)
	}

	return e
}
