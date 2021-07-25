package main

import (
	"encoding/base64"
	"fmt"
	"testing"

	"golang.org/x/crypto/curve25519"
)

func TestKex(t *testing.T) {

	pub, key, err := NewKeypair()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Alice Public key: \t%s\n", b64(pub))
	fmt.Printf("Alice Private key:\t%s\n", b64(key))

	pub2, key2, err := NewKeypair()
	if err != nil {
		panic(err)

	}
	fmt.Printf("Bob Public key: \t%s\n", b64(pub2))
	fmt.Printf("Bob Private key:\t%s\n", b64(key2))

	// exchange public keys here

	out, err := curve25519.X25519(key, pub2)
	if err != nil {
		panic(err)
	}
	out2, err := curve25519.X25519(key2, pub)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Shared key (Alice):\t%s\n", b64(out))
	fmt.Printf("Shared key (Bob):\t%s\n", b64(out2))
}

func b64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
