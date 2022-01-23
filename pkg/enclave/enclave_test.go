package enclave_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"myst/pkg/crypto"
	"myst/pkg/enclave"
)

func TestEnclaveEncryptDecryptLossless(t *testing.T) {
	b := []byte("payload")
	password := []byte("my-password")

	salt, err := crypto.GenerateRandomBytes(uint(crypto.DefaultArgon2IdParams.SaltLength))
	assert.Nil(t, err)

	key := crypto.Argon2Id(password, salt)

	enc, err := enclave.Create(b, key, salt)
	assert.Nil(t, err)

	b2, err := enclave.Unlock(enc, key)
	assert.Nil(t, err)
	assert.Equal(t, b, b2)
}
