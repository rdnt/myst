package crypto_test

import (
	"encoding/base64"
	"testing"

	"golang.org/x/crypto/curve25519"
	"gotest.tools/v3/assert"

	"myst/pkg/crypto"
)

func TestKex(t *testing.T) {
	pub, key, err := crypto.NewCurve25519Keypair()
	assert.NilError(t, err)

	t.Logf("Alice Public key: \t%s\n", b64(pub))
	t.Logf("Alice Private key:\t%s\n", b64(key))

	pub2, key2, err := crypto.NewCurve25519Keypair()
	assert.NilError(t, err)

	t.Logf("Bob Public key: \t%s\n", b64(pub2))
	t.Logf("Bob Private key:\t%s\n", b64(key2))

	// exchange pub keys...

	out, err := curve25519.X25519(key, pub2)
	assert.NilError(t, err)

	out2, err := curve25519.X25519(key2, pub)
	assert.NilError(t, err)

	t.Logf("Shared key (Alice):\t%s\n", b64(out))
	t.Logf("Shared key (Bob):\t%s\n", b64(out2))

	assert.Equal(t, b64(out), b64(out2))
}

func b64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
