package entry

import (
	"myst/pkg/uuid"
)

type Entry struct {
	Id       string
	Website  string
	Username string
	Password string
	Notes    string
}

func (e *Entry) SetWebsite(website string) {
	e.Website = website
}

func (e *Entry) SetUsername(username string) {
	e.Username = username
}

func (e *Entry) SetPassword(password string) {
	e.Password = password
}

func (e *Entry) SetNotes(notes string) {
	e.Notes = notes
}

func New(opts ...Option) Entry {
	e := Entry{
		Id: uuid.New().String(),
	}

	for _, opt := range opts {
		opt(&e)
	}

	return e
}
