package crypto_test

import (
	"crypto/aes"
	"testing"

	"github.com/stretchr/testify/assert"

	"myst/pkg/crypto"
)

func TestPKCS7PadUnpad(t *testing.T) {
	s := ""

	for i := 0; i < 100; i++ {
		s += "a"
		b := []byte(s)

		b2, err := crypto.PKCS7Pad(b, aes.BlockSize)
		assert.Nil(t, err)

		b3, err := crypto.PKCS7Unpad(b2, aes.BlockSize)
		assert.Nil(t, err)

		assert.Equal(t, b, b3)
	}

	//b := []byte("padyloadssssdsww")
	//b2, err := crypto.PKCS7Pad(b, aes.BlockSize)
	//assert.Nil(t, err)
	//b3, err := crypto.PKCS7Unpad(b2, aes.BlockSize)
	//assert.Nil(t, err)
	//assert.Equal(t, b, b3)
}

func FuzzPKCS7PadUnpad(f *testing.F) {
	//f.Add([]byte("payload"))
	//f.Add([]byte("ramdom payload"))
	//f.Add([]byte("nore random payload"))
	//f.Skip(nil)

	f.Fuzz(func(t *testing.T, b []byte) {
		if b == nil || len(b) == 0 {
			t.Skip()
		}
		b2, err := crypto.PKCS7Pad(b, aes.BlockSize)
		if err != nil {
			t.Error(b)
		}
		assert.Nil(t, err)

		b3, err := crypto.PKCS7Unpad(b2, aes.BlockSize)
		assert.Nil(t, err)

		assert.Equal(t, b, b3)
	})
}
