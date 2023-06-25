package enclaverepo

import (
	"crypto/sha256"
	"io/fs"
	"os"

	"github.com/pkg/errors"

	"myst/pkg/crypto"
	"myst/src/client/application"
	"myst/src/client/application/domain/credentials"
	"myst/src/client/application/domain/keystore"
)

type Repository struct {
	path string
}

func New(path string) *Repository {
	r := &Repository{
		path: path,
	}

	return r
}

func (r *Repository) Keystore(key []byte, id string) (keystore.Keystore, error) {
	e, err := r.enclave(key)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get enclave")
	}

	k, err := e.keystore(id)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystore")
	}

	return k, nil
}

func (r *Repository) Keystores(key []byte) (map[string]keystore.Keystore, error) {
	_, err := os.ReadFile(r.enclavePath())
	if errors.Is(err, os.ErrNotExist) {
		return nil, application.ErrInitializationRequired
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to read enclave")
	}

	e, err := r.enclave(key)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get enclave")
	}

	ks, err := e.keystoresWithKeys()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get keystores")
	}

	return ks, nil
}

func (r *Repository) IsInitialized() (bool, error) {
	exists, err := r.enclaveExists()
	if err != nil {
		return false, errors.WithMessage(err, "failed to check if enclave exists")
	}

	if !exists {
		return false, nil
	}

	return true, nil
}

func (r *Repository) enclaveExists() (bool, error) {
	_, err := os.Stat(r.enclavePath())
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, errors.Wrap(err, "failed to stat enclave")
	}

	return true, nil
}

func (r *Repository) DeleteKeystore(key []byte, id string) error {
	e, err := r.enclave(key)
	if err != nil {
		return errors.WithMessage(err, "failed to get enclave")
	}

	err = e.deleteKeystore(id)
	if err != nil {
		return errors.WithMessage(err, "failed to delete keystore")
	}

	err = r.sealAndWrite(key, e)
	if err != nil {
		return errors.WithMessage(err, "failed to seal enclave")
	}

	return nil
}

func (r *Repository) CreateKeystore(key []byte, k keystore.Keystore) (keystore.Keystore, error) {
	e, err := r.enclave(key)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get enclave")
	}

	err = e.addKeystore(k)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to add keystore")
	}

	err = r.sealAndWrite(key, e)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to seal enclave")
	}

	return k, nil
}

// Initialize initializes the enclave with the given password
func (r *Repository) Initialize(password string) ([]byte, error) {
	exists, err := r.enclaveExists()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to check if enclave exists")
	}

	if exists {
		return nil, application.ErrEnclaveExists
	}

	p := crypto.DefaultArgon2IdParams

	salt, err := crypto.GenerateRandomBytes(uint(p.SaltLength))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate salt")
	}

	key := crypto.Argon2Id([]byte(password), salt)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate key")
	}

	publicKey, privateKey, err := crypto.NewCurve25519Keypair()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate keypair")
	}

	e := newEnclave(salt, publicKey, privateKey)

	err = r.sealAndWrite(key, e)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to seal enclave")
	}

	return key, nil
}

func (r *Repository) Authenticate(password string) ([]byte, error) {
	// we don't defer unlock because hash computation is slow

	b, err := os.ReadFile(r.enclavePath())
	if errors.Is(err, fs.ErrNotExist) {
		return nil, application.ErrInitializationRequired
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to read enclave")
	}

	salt := getSaltFromData(b)

	key := crypto.Argon2Id([]byte(password), salt)

	p := crypto.DefaultArgon2IdParams

	mac := b[p.SaltLength : sha256.Size+p.SaltLength]
	b = b[p.SaltLength+sha256.Size:]

	valid := crypto.VerifyHMAC_SHA256(key, mac, b)
	if !valid {
		return nil, application.ErrAuthenticationFailed
	}

	return key, nil
}

func (r *Repository) UpdateKeystore(key []byte, k keystore.Keystore) (keystore.Keystore, error) {
	e, err := r.enclave(key)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get enclave")
	}

	k, err = e.updateKeystore(k)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to update keystore")
	}

	err = r.sealAndWrite(key, e)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to seal enclave")
	}

	return k, nil
}

func (r *Repository) sealAndWrite(key []byte, e *enclave) error {
	b, err := enclaveToJSON(e)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal enclave")
	}

	b, err = crypto.AES256CBC_Encrypt(key, b)
	if err != nil {
		return errors.WithMessage(err, "failed to encrypt enclave")
	}

	// authenticate
	mac := crypto.HMAC_SHA256(key, b)

	// prepend salt and mac to the ciphertext
	b = append(e.salt, append(mac, b...)...)

	err = os.WriteFile(r.enclavePath(), b, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to write enclave")
	}

	return nil
}

func (r *Repository) UpdateCredentials(key []byte, creds credentials.Credentials) (credentials.Credentials, error) {
	e, err := r.enclave(key)
	if err != nil {
		return credentials.Credentials{}, errors.WithMessage(err, "failed to get enclave")
	}

	e.creds = &creds

	err = r.sealAndWrite(key, e)
	if err != nil {
		return credentials.Credentials{}, errors.WithMessage(err, "failed to seal enclave")
	}

	return creds, nil
}

func (r *Repository) Credentials(key []byte) (credentials.Credentials, error) {
	e, err := r.enclave(key)
	if err != nil {
		return credentials.Credentials{}, errors.WithMessage(err, "failed to get enclave")
	}

	rem := e.creds
	if rem == nil {
		return credentials.Credentials{}, application.ErrCredentialsNotFound
	}

	return *rem, nil
}
