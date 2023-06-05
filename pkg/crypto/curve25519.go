package crypto

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/curve25519"
)

func NewCurve25519Keypair() (publicKey []byte, privateKey []byte, err error) {
	var pub [32]byte
	var key [32]byte

	b, err := GenerateRandomBytes(32)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to generate random bytes")
	}

	copy(key[:], b)

	curve25519.ScalarBaseMult(&pub, &key)

	return pub[:], key[:], nil
}
