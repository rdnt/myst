package enclave

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"

	"myst/pkg/crypto"
)

var (
	ErrAuthenticationFailed = fmt.Errorf("authentication failed")
	ErrCorrupted            = fmt.Errorf("corrupted")

	checkKey = []byte{
		0x88, 0x80, 0xc1, 0x7e, 0x69, 0xbe, 0xa3, 0x03,
		0x50, 0x36, 0x5b, 0x52, 0x13, 0x3f, 0x6a, 0xa2,
		0xb2, 0x7a, 0xc0, 0xa2, 0x24, 0x62, 0x53, 0xa3,
		0xf2, 0xa4, 0x93, 0xa2, 0x86, 0x31, 0xb0, 0x69,
	}
)

// Create creates a new encrypted enclave. We use encrypt then MAC so that we
// get integrity of ciphertext. The enclave consists of:
// [salt][mac][ciphertext]
//                 ^
//            [checkKey][payload] (encrypted)
// The sealed enclave must have a length not less than:
// len(salt) + len(mac) + len(encryptedCheckKey)
func Create(b []byte, key []byte, salt []byte) ([]byte, error) {
	// encrypt the check-key along with the payload
	b, err := crypto.AES256CBC_Encrypt(key, append(checkKey, b...))
	if err != nil {
		return nil, err
	}

	// authenticate
	mac := crypto.HMAC_SHA256(key, b)

	// prepend mac to the ciphertext
	b = append(mac, b...)

	// prepend salt to the mac-ciphertext
	b = append(salt, b...)

	return b, nil
}

// Unlock unlocks the enclave, returning the original plaintext (without the
// check-key)
func Unlock(b []byte, key []byte) ([]byte, error) {
	if len(b) == 0 {
		return []byte{}, nil
	}

	p := crypto.DefaultArgon2IdParams

	mac := b[p.SaltLength : sha256.Size+p.SaltLength]
	b = b[p.SaltLength+sha256.Size:]

	valid := crypto.VerifyHMAC_SHA256(key, mac, b)
	if !valid {
		return nil, ErrCorrupted
	}

	// Decrypt keystore
	b, err := crypto.AES256CBC_Decrypt(key, b)
	if err != nil {
		return nil, err
	}

	// verify the decrypted checkKey matches
	if subtle.ConstantTimeCompare(b[:len(checkKey)], checkKey) != 1 {
		return nil, ErrAuthenticationFailed
	}

	// return the decrypted payload (without the checkKey)
	return b[len(checkKey):], nil
}

func GetSaltFromData(b []byte) ([]byte, error) {
	p := crypto.DefaultArgon2IdParams

	if len(b) < int(p.SaltLength) {
		return nil, fmt.Errorf("invalid data")
	}

	return b[:p.SaltLength], nil
}
