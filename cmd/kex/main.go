package main

import (
	"fmt"
	"golang.org/x/crypto/curve25519"

	"myst/pkg/crypto"
)

var eccKeySize = uint(32)

func main() {

	pub, key, err := NewKeypair()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Alice Public key: \t%x\n", pub)
	fmt.Printf("Alice Private key:\t%x\n", key)

	pub2, key2, err := NewKeypair()
	if err != nil {
		panic(err)

	}
	fmt.Printf("Bob Public key: \t%x\n", pub2)
	fmt.Printf("Bob Private key:\t%x\n", key2)

	// exchange public keys here

	out, err := curve25519.X25519(key, pub2)
	if err != nil {
		panic(err)
	}
	out2, err := curve25519.X25519(key2, pub)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Shared key (Alice):\t%x\n", out)
	fmt.Printf("Shared key (Bob):\t%x\n", out2)
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
