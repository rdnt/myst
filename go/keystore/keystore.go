package keystore

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/sht/myst/go/crypto"
	"io"
	"os"
)

// store holds the keystore in-memory
var store Keystore
var key []byte

var ErrNoEntry = fmt.Errorf("entry not found")
var ErrNotModified = fmt.Errorf("keystore not modified")
var ErrAuthFailed = fmt.Errorf("authentication failed")

// Keystore holds all the data
type Keystore struct {
	modified bool
	Entries  map[string]*Entry `json:"entries"`
}

// Entry represents a website and its saved data
type Entry struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Type     string `json:"type"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Decrypt(data []byte, password string) (*Keystore, error) {
	if len(data) == 0 {
		store := Keystore{
			modified: true,
			Entries:  map[string]*Entry{},
		}
		return &store, nil
	}

	p := crypto.GetArgon2IdParams()

	salt := data[:p.SaltLength]
	mac := data[p.SaltLength : sha256.Size+p.SaltLength]
	data = data[p.SaltLength+sha256.Size:]

	key := crypto.Argon2Id([]byte(password), salt)

	valid := crypto.VerifyHMAC_SHA256(key, mac, data)
	if !valid {
		return nil, ErrAuthFailed
	}

	// Decrypt keystore
	data, err := crypto.AES256CBC_Decrypt(key, data)
	if err != nil {
		return nil, err
	}
	// Decode decrypted keystore from json
	err = json.Unmarshal(data, &store)
	if err != nil {
		return nil, err
	}
	// Return the keystore
	return &store, nil
}

func Load(path string) ([]byte, error) {
	// Acquire file handle
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	// Read the keystore file into a bytes array
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, f)
	if err != nil {
		return nil, err
	}

	data := buf.Bytes()

	return data, nil
}

func (ks *Keystore) Save(path, password string) error {
	if !ks.modified {
		return ErrNotModified
	}
	p := crypto.GetArgon2IdParams()
	salt, err := crypto.GenerateRandomBytes(uint(p.SaltLength))
	if err != nil {
		return err
	}

	key := crypto.Argon2Id([]byte(password), salt)

	// Acquire file handle
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	// Encode to json
	data, err := json.Marshal(ks)
	if err != nil {
		return err
	}

	// Encrypt keystore
	data, err = crypto.AES256CBC_Encrypt(key, data)
	if err != nil {
		return err
	}

	// Authenticate ciphertext
	mac := crypto.HMAC_SHA256(key, data)
	data = append(mac, data...)
	data = append(salt, data...)

	// Write encrypted keystore to file
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// String implements the stringer interface
func (ks *Keystore) String() string {
	data, err := json.Marshal(ks)
	if err != nil {
		return ""
	}
	return string(data)
}

func (ks *Keystore) Add(email, password string) (*Entry, error) {
	if ks.Entries == nil {
		ks.Entries = map[string]*Entry{}
	}
	var entry *Entry
	for {
		uid, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}
		id := uid.String()
		if _, exists := ks.Entries[id]; !exists {
			entry = &Entry{
				ID:       id,
				Email:    email,
				Password: password,
			}
			ks.Entries[id] = entry
			ks.modified = true
			break
		}
	}
	return entry, nil
}

func (ks *Keystore) Remove(id string) (bool, error) {
	if ks.Entries == nil {
		return false, ErrNoEntry
	}
	if _, exists := ks.Entries[id]; !exists {
		return false, ErrNoEntry
	}
	delete(ks.Entries, id)
	ks.modified = true
	return true, nil
}

func (ks *Keystore) Get(id string) (*Entry, error) {
	if ks.Entries == nil {
		return nil, ErrNoEntry
	}
	if _, exists := ks.Entries[id]; !exists {
		return nil, ErrNoEntry
	}
	return ks.Entries[id], nil

}
