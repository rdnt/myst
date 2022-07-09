package main

import (
	"golang.org/x/crypto/curve25519"

	"myst/pkg/crypto"
)

var eccKeySize = uint(32)

func main() {

}

func NewKeypair() ([]byte, []byte, error) {
	var pub [32]byte
	var key [32]byte

	b, err := crypto.GenerateRandomBytes(eccKeySize)
	if err != nil {
		return nil, nil, err
	}
	copy(key[:], b)

	curve25519.ScalarBaseMult(&pub, &key)

	return pub[:], key[:], nil
}
