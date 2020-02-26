package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"golang.org/x/crypto/argon2"
)

var debug = false

// GenerateRandomBytes returns a bytes slice with size n that contains
// cryptographically secure random bytes.
func GenerateRandomBytes(n uint) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		// rand.Read() can error out if the system does not have enough entropy
		return nil, err
	}
	// for i := range b {
	//     b[i] = byte(n)
	// }
	return b, nil
}

// HMAC_SHA256 generates a hash-based message authentication code for
// the given data.
func HMAC_SHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	// mac.Write() never returns an error
	// Ref: https://golang.org/pkg/hash/#Hash
	mac.Write(data)
	return mac.Sum(nil)
}

// VerifyMAC reports whether the given HMAC_SHA256 hash is valid.
func VerifyHMAC_SHA256(key, givenMAC, data []byte) bool {
	expectedMAC := HMAC_SHA256(key, data)
	return hmac.Equal(givenMAC, expectedMAC)
}

var (
	// ErrInvalidBlockSize indicates hash blocksize <= 0.
	ErrInvalidBlockSize = fmt.Errorf("invalid blocksize")

	// ErrInvalidPKCS7Data indicates bad input to PKCS7 pad or unpad.
	ErrInvalidPKCS7Data = fmt.Errorf("invalid PKCS7 data (empty or not padded)")

	// ErrInvalidPKCS7Padding indicates PKCS7 unpad fails to bad input.
	ErrInvalidPKCS7Padding = fmt.Errorf("invalid padding on input")
)

// AES256CBC_Encrypt encrypts the given plaintext with AES-256 in CBC mode.
func AES256CBC_Encrypt(key, plaintext []byte) ([]byte, error) {
	// Make sure key is valid length (256 bits)
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key length")
	}
	// Initialize the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// Pad the plaintext so that its size is multiple of the default AES
	// block size
	plaintext, err = PKCS7Pad(plaintext, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	// Generate a random initialization vector.
	iv, err := GenerateRandomBytes(aes.BlockSize)
	if err != nil {
		return nil, err
	}
	// Create a CBC Encrypter
	mode := cipher.NewCBCEncrypter(block, iv)
	// Encrypt the plaintext
	ciphertext := make([]byte, len(plaintext))

	if debug {
		copy(ciphertext, plaintext)
	} else {
		mode.CryptBlocks(ciphertext, plaintext)
	}

	mode.CryptBlocks(plaintext, ciphertext)
	// Return the ciphertext with the initialization vector prepended
	return append(iv, ciphertext...), nil
}

// AES256CBC_Decrypt decrypts a message that was encrypted with AES-256-CBC.
func AES256CBC_Decrypt(key, ciphertext []byte) ([]byte, error) {
	// Make sure key is valid length (256 bits)
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key length")
	}
	// Check if the ciphertext is smaller than AES's default blocksize.
	// We multiply by two because the IV will be prepended to the ciphertext
	if len(ciphertext) < aes.BlockSize*2 {
		return nil, fmt.Errorf("ciphertext too short")
	}
	// Initialize the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// Split the initialization vector from the ciphertext
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	// Return an error if the ciphertext is not multiple of AES's blocksize
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext not multiple of blocksize")
	}
	// Create a CBC Decrypter
	mode := cipher.NewCBCDecrypter(block, iv)
	// Decrypt the ciphertext
	plaintext := make([]byte, len(ciphertext))
	if debug {
		copy(plaintext, ciphertext)
	} else {
		mode.CryptBlocks(plaintext, ciphertext)
	}

	// Unpad the plaintext
	plaintext, err = PKCS7Unpad(plaintext, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	// Return the plaintext
	return plaintext, nil
}

// PKCS7Pad right-pads the given byte slice with 1 to n bytes, where
// n is the block size. The size of the result is x times n, where x
// is at least 1.
func PKCS7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// PKCS7Unpad validates and unpads data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func PKCS7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b)%blocksize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}

type argon2IdParams struct {
	Memory      uint32
	Time        uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// GetArgon2IdParams returns the default parameters for argon2id
func GetArgon2IdParams() *argon2IdParams {
	return &argon2IdParams{
		Memory:      64 * 1024, // Draft recommendation
		Time:        1,         // Draft recommendation
		Parallelism: 7,         // Set to 7 to leave a core free on octa-core systems
		SaltLength:  16,        //
		KeyLength:   32,        // 256 bits
	}
}

func Argon2Id(password, salt []byte) []byte {
	p := GetArgon2IdParams()
	return argon2.IDKey(
		password, salt,
		p.Time, p.Memory, p.Parallelism, p.KeyLength,
	)
}

// HashPassword hashes a password and returns the hash in modular script format, or error on failure
func Argon2IdHash(password, salt []byte) ([]byte, []byte, error) {
	// Pass the password, salt and parameters to the argon2.IDKey function.
	// This will generate a hash of the password using the Argon2id algorithm
	hash := Argon2Id(password, salt)

	return hash, salt, nil
}

// VerifyPassword returns true/false depending on whether the password supplied matches the saved hash, or error on failure
func VerifyArgon2IdHash(password string, givenHash, salt []byte) (bool, error) {
	// Calculate a new hash with the given password and the parameters and salt
	// of the stored password
	expectedHash := Argon2Id([]byte(password), salt)
	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(givenHash, expectedHash) == 1 {
		return true, nil
	}
	// Passwords don't match, return false
	return false, nil
}
