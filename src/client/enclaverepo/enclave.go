package enclaverepo

import (
	"crypto/sha256"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"

	"myst/pkg/crypto"
	"myst/src/client/application/domain/credentials"
	"myst/src/client/application/domain/keystore"
	"myst/src/server/application"
)

type enclave struct {
	encSalt   []byte
	signSalt  []byte
	keystores map[string]keystore.Keystore
	keys      map[string][]byte
	creds     credentials.Credentials
}

func newEnclave(encSalt, signSalt, publicKey, privateKey []byte) *enclave {
	return &enclave{
		keystores: map[string]keystore.Keystore{},
		keys:      map[string][]byte{},
		encSalt:   encSalt,
		signSalt:  signSalt,
		creds: credentials.Credentials{
			PublicKey:  publicKey,
			PrivateKey: privateKey,
		},
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

func (e *enclave) updateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	k.UpdatedAt = time.Now()
	k.Version++

	e.keystores[k.Id] = k

	return k, nil
}

func (e *enclave) deleteKeystore(id string) error {
	delete(e.keystores, id)
	delete(e.keys, id)

	return nil
}

func (r *Repository) enclavePath() string {
	return path.Join(r.path, "data.myst")
}

func (r *Repository) enclave(keypair []byte) (*enclave, error) {
	b, err := os.ReadFile(r.enclavePath())
	if err != nil {
		return nil, errors.Wrap(err, "failed to read enclave file")
	}

	p := crypto.DefaultArgon2IdParams

	encKey := keypair[0:p.KeyLength]
	signKey := keypair[p.KeyLength : p.KeyLength*2]

	encSalt, signSalt := getSaltsFromData(b)

	mac := b[p.SaltLength*2 : sha256.Size+p.SaltLength*2]
	b = b[p.SaltLength*2+sha256.Size:]

	valid := crypto.VerifyHMAC_SHA256(signKey, mac, b)
	if !valid {
		return nil, application.ErrAuthenticationFailed
	}

	b, err = crypto.AES256CBC_Decrypt(encKey, b)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to decrypt enclave")
	}

	encJson, err := enclaveFromJSON(b, encSalt, signSalt)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create enclave from json")
	}

	return encJson, nil
}

func getSaltsFromData(b []byte) (encSalt, signSalt []byte) {
	encSaltLen := crypto.DefaultArgon2IdParams.SaltLength
	signSaltLen := crypto.DefaultArgon2IdParams.SaltLength

	return b[:encSaltLen], b[encSaltLen : encSaltLen+signSaltLen]
}
