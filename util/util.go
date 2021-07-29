package util

import (
	"github.com/lithammer/shortuuid/v3"
)

// NewUUID returns a new pretty uuid
func NewUUID() string {
	return shortuuid.New()
}
