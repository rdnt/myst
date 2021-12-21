package enclave

import (
	"crypto/sha256"
	"fmt"

	"myst/pkg/crypto"
)

var ErrAuthenticationFailed = fmt.Errorf("authentication failed")

func Encrypt(b []byte, key []byte, salt []byte) ([]byte, error) {
	// Encrypt keystore
	b, err := crypto.AES256CBC_Encrypt(key, b)
	if err != nil {
		return nil, err
	}

	// Authenticate ciphertext
	mac := crypto.HMAC_SHA256(key, b)
	b = append(mac, b...)
	b = append(salt, b...)

	// Return encrypted keystore as bytes slice
	return b, nil
}

func Decrypt(b []byte, key []byte) ([]byte, error) {
	if len(b) == 0 {
		return []byte{}, nil
	}

	p := crypto.DefaultArgon2IdParams

	mac := b[p.SaltLength : sha256.Size+p.SaltLength]
	b = b[p.SaltLength+sha256.Size:]

	valid := crypto.VerifyHMAC_SHA256(key, mac, b)
	if !valid {
		return nil, ErrAuthenticationFailed
	}

	// Decrypt keystore
	return crypto.AES256CBC_Decrypt(key, b)
}

func GetSaltFromData(data []byte) ([]byte, error) {
	p := crypto.DefaultArgon2IdParams

	if len(data) < int(p.SaltLength) {
		return nil, fmt.Errorf("invalid data")
	}
	salt := data[:p.SaltLength]
	return salt, nil
}