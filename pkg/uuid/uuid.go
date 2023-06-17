package uuid

import "github.com/lithammer/shortuuid/v3"

type UUID string

func New() UUID {
	return UUID(shortuuid.New())
}

func (u UUID) String() string {
	return string(u)
}
