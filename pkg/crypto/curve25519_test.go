package crypto_test

import (
	"encoding/base64"
	"testing"

	"golang.org/x/crypto/curve25519"
	"gotest.tools/v3/assert"

	"myst/pkg/crypto"
)

func TestKeyExchange(t *testing.T) {
	pub, key, err := crypto.NewCurve25519Keypair()
	assert.NilError(t, err)

	pub2, key2, err := crypto.NewCurve25519Keypair()
	assert.NilError(t, err)

	// exchange pub keys...

	out, err := curve25519.X25519(key, pub2)
	assert.NilError(t, err)

	out2, err := curve25519.X25519(key2, pub)
	assert.NilError(t, err)

	assert.Equal(t, b64(out), b64(out2))
}

func b64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
