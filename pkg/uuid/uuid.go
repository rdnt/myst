package uuid

import "github.com/lithammer/shortuuid/v3"

type UUID string

func (u UUID) String() string {
	return string(u)
}

func New() UUID {
	return UUID(shortuuid.New())
}
