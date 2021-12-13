package keystorerepo

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"myst/pkg/crypto"
)

func Encrypt(b []byte, key []byte) ([]byte, error) {
	// Encode to json
	b, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	p := crypto.DefaultArgon2IdParams

	salt, err := crypto.GenerateRandomBytes(uint(p.SaltLength))
	if err != nil {
		return nil, err
	}

	// Encrypt keystore
	b, err = crypto.AES256CBC_Encrypt(key, b)
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
		return nil, fmt.Errorf("authentication failed")
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
