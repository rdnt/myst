package crypto_test

import (
	"crypto/aes"
	"testing"

	"gotest.tools/v3/assert"

	"myst/pkg/crypto"
)

func TestPKCS7PadUnpad(t *testing.T) {
	s := ""

	for i := 0; i < 100; i++ {
		s += "a"
		b := []byte(s)

		b2, err := crypto.PKCS7Pad(b, aes.BlockSize)
		assert.NilError(t, err)

		b3, err := crypto.PKCS7Unpad(b2, aes.BlockSize)
		assert.NilError(t, err)

		assert.DeepEqual(t, b, b3)
	}
}

func FuzzPKCS7PadUnpad(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		if len(b) == 0 {
			t.Skip()
		}
		b2, err := crypto.PKCS7Pad(b, aes.BlockSize)
		if err != nil {
			t.Error(b)
		}
		assert.NilError(t, err)

		b3, err := crypto.PKCS7Unpad(b2, aes.BlockSize)
		assert.NilError(t, err)

		assert.Equal(t, b, b3)
	})
}
