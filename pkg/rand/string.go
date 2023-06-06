package rand

import (
	"crypto/rand"

	"github.com/pkg/errors"
)

func String(length int) (string, error) {
	result := make([]byte, length)

	_, err := rand.Read(result)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate random string")
	}

	for i := 0; i < length; i++ {
		result[i] &= 0x7F

		for result[i] < 32 || result[i] == 127 {
			_, err = rand.Read(result[i : i+1])
			if err != nil {
				return "", errors.Wrap(err, "failed to generate random string")
			}

			result[i] &= 0x7F
		}
	}

	return string(result), nil
}

// Bytes returns a bytes slice with size n that contains
// cryptographically secure random bytes.
func Bytes(n uint) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate random bytes")
	}

	return b, nil
}
