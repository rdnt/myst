package enclaverepo

import (
	"crypto/sha256"
	"github.com/pkg/errors"
	"myst/src/server/application"
	"os"
	"path"
	"time"

	"myst/pkg/crypto"
	"myst/src/client/application/domain/credentials"
	"myst/src/client/application/domain/keystore"
)

type enclave struct {
	salt      []byte
	keystores map[string]keystore.Keystore
	keys      map[string][]byte
	creds     *credentials.Credentials
}

func newEnclave(salt []byte) *enclave {
	return &enclave{
		keystores: map[string]keystore.Keystore{},
		keys:      map[string][]byte{},
		salt:      salt,
	}
}

func (e *enclave) addKeystore(k keystore.Keystore) error {
	p := crypto.DefaultArgon2IdParams

	key, err := crypto.GenerateRandomBytes(uint(p.KeyLength))
	if err != nil {
		return errors.WithMessage(err, "failed to generate key")
	}

	k.CreatedAt = time.Now()
	k.UpdatedAt = time.Now()

	e.keystores[k.Id] = k
	e.keys[k.Id] = key

	return nil
}

func (e *enclave) keystoresWithKeys() (map[string]keystore.Keystore, error) {
	ks := map[string]keystore.Keystore{}

	for _, k := range e.keystores {
		keystoreKey, ok := e.keys[k.Id]
		if !ok {
			return nil, application.ErrKeystoreNotFound
		}

		k.Key = keystoreKey
		ks[k.Id] = k
	}

	return ks, nil
}

func (e *enclave) keystore(id string) (keystore.Keystore, error) {
	k, ok := e.keystores[id]
	if !ok {
		return keystore.Keystore{}, application.ErrKeystoreNotFound
	}

	keystoreKey, ok := e.keys[id]
	if !ok {
		return keystore.Keystore{}, application.ErrKeystoreNotFound
	}

	k.Key = keystoreKey

	return k, nil
}

func (e *enclave) updateKeystore(k keystore.Keystore) error {
	k.UpdatedAt = time.Now()
	k.Version++

	e.keystores[k.Id] = k

	return nil
}

func (e *enclave) deleteKeystore(id string) error {
	delete(e.keystores, id)
	delete(e.keys, id)

	return nil
}

func (r *Repository) enclavePath() string {
	return path.Join(r.path, "data.myst")
}

func (r *Repository) enclave(argon2idKey []byte) (*enclave, error) {
	b, err := os.ReadFile(r.enclavePath())
	if err != nil {
		return nil, errors.Wrap(err, "failed to read enclave file")
	}

	salt := getSaltFromData(b)

	p := crypto.DefaultArgon2IdParams

	mac := b[p.SaltLength : sha256.Size+p.SaltLength]
	b = b[p.SaltLength+sha256.Size:]

	valid := crypto.VerifyHMAC_SHA256(argon2idKey, mac, b)
	if !valid {
		return nil, application.ErrAuthenticationFailed
	}

	b, err = crypto.AES256CBC_Decrypt(argon2idKey, b)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to decrypt enclave")
	}

	encJson, err := enclaveFromJSON(b, salt)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create enclave from json")
	}

	return encJson, nil
}

func getSaltFromData(b []byte) []byte {
	return b[:crypto.DefaultArgon2IdParams.SaltLength]
}
