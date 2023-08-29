package enclaverepo

import (
	"crypto/sha256"
	"io/fs"
	"os"
	"sync"

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

func (r *Repository) Keystore(keypair []byte, id string) (keystore.Keystore, error) {
	e, err := r.enclave(keypair)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get enclave")
	}

	k, err := e.keystore(id)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystore")
	}

	return k, nil
}

func (r *Repository) Keystores(keypair []byte) (map[string]keystore.Keystore, error) {
	_, err := os.ReadFile(r.enclavePath())
	if errors.Is(err, os.ErrNotExist) {
		return nil, application.ErrInitializationRequired
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to read enclave")
	}

	e, err := r.enclave(keypair)
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

func (r *Repository) DeleteKeystore(keypair []byte, id string) error {
	e, err := r.enclave(keypair)
	if err != nil {
		return errors.WithMessage(err, "failed to get enclave")
	}

	err = e.deleteKeystore(id)
	if err != nil {
		return errors.WithMessage(err, "failed to delete keystore")
	}

	err = r.sealAndWrite(keypair, e)
	if err != nil {
		return errors.WithMessage(err, "failed to seal enclave")
	}

	return nil
}

func (r *Repository) CreateKeystore(keypair []byte, k keystore.Keystore) (keystore.Keystore, error) {
	e, err := r.enclave(keypair)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get enclave")
	}

	err = e.addKeystore(k)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to add keystore")
	}

	err = r.sealAndWrite(keypair, e)
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

	encSalt, err := crypto.GenerateRandomBytes(uint(p.SaltLength))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate salt")
	}

	signSalt, err := crypto.GenerateRandomBytes(uint(p.SaltLength))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate salt")
	}

	var encKey, signKey []byte

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		encKey = crypto.Argon2Id([]byte(password), encSalt)
	}()

	go func() {
		defer wg.Done()

		signKey = crypto.Argon2Id([]byte(password), signSalt)
	}()

	wg.Wait()

	publicKey, privateKey, err := crypto.NewCurve25519Keypair()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to generate keypair")
	}

	e := newEnclave(encSalt, signSalt, publicKey, privateKey)

	keypair := append(encKey, signKey...)

	err = r.sealAndWrite(keypair, e)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to seal enclave")
	}

	return keypair, nil
}

func (r *Repository) Authenticate(password string) (keypair []byte, err error) {
	// we don't defer unlock because hash computation is slow

	b, err := os.ReadFile(r.enclavePath())
	if errors.Is(err, fs.ErrNotExist) {
		return nil, application.ErrInitializationRequired
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to read enclave")
	}

	encSalt, signSalt := getSaltsFromData(b)

	var encKey, signKey []byte

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		encKey = crypto.Argon2Id([]byte(password), encSalt)
	}()

	go func() {
		defer wg.Done()

		signKey = crypto.Argon2Id([]byte(password), signSalt)
	}()

	wg.Wait()

	p := crypto.DefaultArgon2IdParams

	mac := b[p.SaltLength*2 : sha256.Size+p.SaltLength*2]
	b = b[p.SaltLength*2+sha256.Size:]

	valid := crypto.VerifyHMAC_SHA256(signKey, mac, b)
	if !valid {
		return nil, application.ErrAuthenticationFailed
	}

	return append(encKey, signKey...), nil
}

func (r *Repository) UpdateKeystore(keypair []byte, k keystore.Keystore) (keystore.Keystore, error) {
	e, err := r.enclave(keypair)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get enclave")
	}

	k, err = e.updateKeystore(k)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to update keystore")
	}

	err = r.sealAndWrite(keypair, e)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to seal enclave")
	}

	return k, nil
}

func (r *Repository) sealAndWrite(keypair []byte, e *enclave) error {
	b, err := enclaveToJSON(e)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal enclave")
	}

	p := crypto.DefaultArgon2IdParams

	encKey := keypair[:p.KeyLength]
	signKey := keypair[p.KeyLength:]

	b, err = crypto.AES256CBC_Encrypt(encKey, b)
	if err != nil {
		return errors.WithMessage(err, "failed to encrypt enclave")
	}

	// authenticate
	mac := crypto.HMAC_SHA256(signKey, b)

	// prepend salt and mac to the ciphertext
	b = append(e.encSalt, append(e.signSalt, append(mac, b...)...)...)

	err = os.WriteFile(r.enclavePath(), b, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to write enclave")
	}

	return nil
}

func (r *Repository) UpdateCredentials(keypair []byte, creds credentials.Credentials) (credentials.Credentials, error) {
	e, err := r.enclave(keypair)
	if err != nil {
		return credentials.Credentials{}, errors.WithMessage(err, "failed to get enclave")
	}

	e.creds = creds

	err = r.sealAndWrite(keypair, e)
	if err != nil {
		return credentials.Credentials{}, errors.WithMessage(err, "failed to seal enclave")
	}

	return creds, nil
}

func (r *Repository) Credentials(keypair []byte) (credentials.Credentials, error) {
	e, err := r.enclave(keypair)
	if err != nil {
		return credentials.Credentials{}, errors.WithMessage(err, "failed to get enclave")
	}

	if e.creds.Address == "" {
		return credentials.Credentials{}, application.ErrCredentialsNotFound
	}

	return e.creds, nil
}
