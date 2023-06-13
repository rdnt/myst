package enclaverepo

import (
	"crypto/sha256"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"

	"myst/pkg/crypto"
	"myst/src/client/application"
	"myst/src/client/application/domain/credentials"
	"myst/src/client/application/domain/keystore"
)

type Repository struct {
	mux             sync.Mutex
	key             []byte
	lastHealthCheck time.Time
	path            string
}

func New(path string) *Repository {
	r := &Repository{
		path: path,
	}

	go r.startHealthCheck()

	return r
}

func (r *Repository) HealthCheck() {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.lastHealthCheck = time.Now()
}

func (r *Repository) Keystore(id string) (keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return keystore.Keystore{}, application.ErrAuthenticationRequired
	}

	e, err := r.enclave(r.key)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get enclave")
	}

	k, err := e.keystore(id)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get keystore")
	}

	return k, nil
}

func (r *Repository) Keystores() (map[string]keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, err := os.ReadFile(r.enclavePath())
	if errors.Is(err, os.ErrNotExist) {
		return nil, application.ErrInitializationRequired
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to read enclave")
	}

	if r.key == nil {
		return nil, application.ErrAuthenticationRequired
	}

	e, err := r.enclave(r.key)
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
	r.mux.Lock()
	defer r.mux.Unlock()

	exists, err := r.enclaveExists()
	if err != nil {
		return false, errors.WithMessage(err, "failed to check if enclave exists")
	}

	if !exists {
		return false, nil
	}

	if r.key == nil {
		return false, application.ErrAuthenticationRequired
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

func (r *Repository) DeleteKeystore(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return application.ErrAuthenticationRequired
	}

	e, err := r.enclave(r.key)
	if err != nil {
		return errors.WithMessage(err, "failed to get enclave")
	}

	err = e.deleteKeystore(id)
	if err != nil {
		return errors.WithMessage(err, "failed to delete keystore")
	}

	err = r.sealAndWrite(e)
	if err != nil {
		return errors.WithMessage(err, "failed to seal enclave")
	}

	return nil
}

func (r *Repository) CreateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return keystore.Keystore{}, application.ErrAuthenticationRequired
	}

	e, err := r.enclave(r.key)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get enclave")
	}

	err = e.addKeystore(k)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to add keystore")
	}

	err = r.sealAndWrite(e)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to seal enclave")
	}

	return k, nil
}

// Initialize initializes the enclave with the given password
func (r *Repository) Initialize(password string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	exists, err := r.enclaveExists()
	if err != nil {
		return errors.WithMessage(err, "failed to check if enclave exists")
	}

	if exists {
		return errors.New("enclave already initialized")
	}

	p := crypto.DefaultArgon2IdParams

	salt, err := crypto.GenerateRandomBytes(uint(p.SaltLength))
	if err != nil {
		return errors.WithMessage(err, "failed to generate salt")
	}

	key := crypto.Argon2Id([]byte(password), salt)
	if err != nil {
		return errors.WithMessage(err, "failed to generate key")
	}

	e := newEnclave(salt)

	r.key = key

	err = r.sealAndWrite(e)
	if err != nil {
		return errors.WithMessage(err, "failed to seal enclave")
	}

	return nil
}

func (r *Repository) Authenticate(password string) error {
	r.mux.Lock()
	// we don't unlock right away because argon2id computation is (purposefully) slow

	b, err := os.ReadFile(r.enclavePath())
	if err != nil {
		r.mux.Unlock()
		return errors.Wrap(err, "failed to read enclave")
	}

	r.mux.Unlock()

	salt := getSaltFromData(b)

	key := crypto.Argon2Id([]byte(password), salt)

	p := crypto.DefaultArgon2IdParams

	mac := b[p.SaltLength : sha256.Size+p.SaltLength]
	b = b[p.SaltLength+sha256.Size:]

	valid := crypto.VerifyHMAC_SHA256(key, mac, b)
	if !valid {
		return application.ErrAuthenticationFailed
	}

	r.mux.Lock()

	r.key = key

	r.mux.Unlock()

	return nil
}

func (r *Repository) UpdateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return keystore.Keystore{}, application.ErrAuthenticationRequired
	}

	e, err := r.enclave(r.key)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get enclave")
	}

	err = e.updateKeystore(k)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to update keystore")
	}

	err = r.sealAndWrite(e)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to seal enclave")
	}

	return k, nil
}

func (r *Repository) startHealthCheck() {
	ticker := time.NewTicker(20 * time.Second) // TODO @rdnt @@@: make this configurable
	defer ticker.Stop()

	for range ticker.C {
		r.mux.Lock()
		elapsed := time.Since(r.lastHealthCheck)

		if elapsed < time.Minute {
			r.mux.Unlock()

			continue
		}

		if r.key != nil {
			fmt.Println("Health check failed")
			// TODO: figure out why healthcheck causes Disconnected
			//  and can't be restored; @@@ @rdnt is this still reproducible?
			r.key = nil
		}

		r.mux.Unlock()
	}
}

func (r *Repository) sealAndWrite(e *enclave) error {
	b, err := enclaveToJSON(e)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal enclave")
	}

	b, err = crypto.AES256CBC_Encrypt(r.key, b)
	if err != nil {
		return errors.WithMessage(err, "failed to encrypt enclave")
	}

	// authenticate
	mac := crypto.HMAC_SHA256(r.key, b)

	// prepend salt and mac to the ciphertext
	b = append(e.salt, append(mac, b...)...)

	err = os.WriteFile(r.enclavePath(), b, 0600)
	if err != nil {
		return errors.Wrap(err, "failed to write enclave")
	}

	return nil
}

func (r *Repository) UpdateCredentials(creds credentials.Credentials) (credentials.Credentials, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return credentials.Credentials{}, application.ErrAuthenticationRequired
	}

	e, err := r.enclave(r.key)
	if err != nil {
		return credentials.Credentials{}, errors.WithMessage(err, "failed to get enclave")
	}

	e.creds = &creds

	err = r.sealAndWrite(e)
	if err != nil {
		return credentials.Credentials{}, errors.WithMessage(err, "failed to seal enclave")
	}

	return creds, nil
}

func (r *Repository) Credentials() (credentials.Credentials, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return credentials.Credentials{}, application.ErrAuthenticationRequired
	}

	e, err := r.enclave(r.key)
	if err != nil {
		return credentials.Credentials{}, errors.WithMessage(err, "failed to get enclave")
	}

	rem := e.creds
	if rem == nil {
		return credentials.Credentials{}, application.ErrCredentialsNotFound
	}

	return *rem, nil
}
