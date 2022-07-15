package enclave_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"myst/pkg/crypto"
	"myst/src/client/pkg/enclave"
)

func TestEnclaveEncryptDecryptLossless(t *testing.T) {
	b := []byte("payload")
	password := []byte("my-password")

	salt, err := crypto.GenerateRandomBytes(uint(crypto.DefaultArgon2IdParams.SaltLength))
	assert.NilError(t, err)

	key := crypto.Argon2Id(password, salt)

	enc, err := enclave.Create(b, key, salt)
	assert.NilError(t, err)

	b2, err := enclave.Unlock(enc, key)
	assert.NilError(t, err)
	assert.DeepEqual(t, b, b2)
}
