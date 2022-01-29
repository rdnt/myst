package entry

import (
	"errors"

	"myst/pkg/uuid"
)

var (
	ErrInvalidLabel    = errors.New("invalid label")
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")
)

type Entry struct {
	id       string
	label    string
	username string
	password string
}

func (e *Entry) Id() string {
	return e.id
}

func (e *Entry) Label() string {
	return e.label
}

func (e *Entry) SetLabel(label string) {
	e.label = label
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

func New(opts ...Option) Entry {
	e := Entry{
		id: uuid.New().String(),
	}

	for _, opt := range opts {
		opt(&e)
	}

	return e
}
