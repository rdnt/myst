package keystorerepo

import (
	"crypto/sha256"
	"fmt"
	"myst/internal/client/application/domain/enclave"
	"myst/internal/client/application/keystoreservice"
	"myst/internal/client/enclaverepo"
	"myst/pkg/logger"
	"os"
	"sync"
	"time"

	"myst/internal/client/application/domain/keystore"
	"myst/pkg/crypto"

	"github.com/pkg/errors"
)

var log = logger.New("keystorerepo", logger.Green)

type Repository struct {
	mux             sync.Mutex
	path            string
	key             []byte
	lastHealthCheck time.Time
	enclaves        *enclaverepo.Repository
}

func New(path string) (*Repository, error) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, err
	}

	r := &Repository{
		path: path,
	}

	go r.startHealthCheck()

	r.enclaves, err = enclaverepo.New("./data")
	if err != nil {
		return nil, err
	}

	// TODO: remove simulated health check
	//go func() {
	//	for {
	//		time.Sleep(10 * time.Second)
	//		r.HealthCheck()
	//	}
	//}()

	return r, nil
}

func (r *Repository) HealthCheck() {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.lastHealthCheck = time.Now()
}

func (r *Repository) Keystore(id string) (*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return nil, fmt.Errorf("authentication required")
	}

	e, err := r.enclaves.Enclave(r.key)
	if err != nil {
		return nil, err
	}

	k, err := e.Keystore(id)
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (r *Repository) Keystores() (map[string]*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, err := os.ReadFile(r.enclaves.EnclavePath())
	if errors.Is(err, os.ErrNotExist) {
		return nil, keystoreservice.ErrInitializationRequired
	} else if err != nil {
		return nil, err
	}

	if r.key == nil {
		return nil, keystoreservice.ErrAuthenticationRequired
	}

	e, err := r.enclaves.Enclave(r.key)
	if err != nil {
		return nil, err
	}

	return e.Keystores(), nil
}

func (r *Repository) Delete(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return fmt.Errorf("authentication required")
	}

	e, err := r.enclaves.Enclave(r.key)
	if err != nil {
		return err
	}

	err = e.DeleteKeystore(id)
	if err != nil {
		return err
	}

	b, err := enclaverepo.MarshalEnclave(e)
	if err != nil {
		return err
	}

	b, err = r.seal(b, r.key, e.Salt())
	if err != nil {
		return err
	}

	err = r.enclaves.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Create(opts ...keystore.Option) (*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return nil, fmt.Errorf("authentication required")
	}

	k := keystore.New(opts...)

	e, err := r.enclaves.Enclave(r.key)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get enclave")
	}

	err = e.AddKeystore(k)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to add keystore")
	}

	b, err := enclaverepo.MarshalEnclave(e)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to marshal enclave")
	}

	b, err = r.seal(b, r.key, e.Salt())
	if err != nil {
		return nil, errors.WithMessage(err, "failed to seal enclave")
	}

	err = r.enclaves.Write(b)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to write enclave")
	}

	return k, nil
}

// Initialize initializes the enclave with the given password
func (r *Repository) Initialize(password string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	// TODO: check if enclave is already created, if it is, return error
	e, err := enclave.New()
	if err != nil {
		return err
	}

	b, err := enclaverepo.MarshalEnclave(e)
	if err != nil {
		return err
	}

	p := crypto.DefaultArgon2IdParams

	salt, err := crypto.GenerateRandomBytes(uint(p.SaltLength))
	if err != nil {
		return err
	}

	key := crypto.Argon2Id([]byte(password), salt)
	if err != nil {
		return err
	}

	b, err = r.seal(b, key, salt)
	if err != nil {
		return err
	}

	err = r.enclaves.Write(b)
	if err != nil {
		return err
	}

	r.key = key

	return nil
}

func (r *Repository) Authenticate(password string) error {
	r.mux.Lock()

	b, err := os.ReadFile(r.enclaves.EnclavePath())
	if err != nil {
		r.mux.Unlock()
		return err
	}

	r.mux.Unlock()

	salt, err := enclaverepo.GetSaltFromData(b)
	if err != nil {
		return err
	}

	key := crypto.Argon2Id([]byte(password), salt)

	p := crypto.DefaultArgon2IdParams

	mac := b[p.SaltLength : sha256.Size+p.SaltLength]
	b = b[p.SaltLength+sha256.Size:]

	valid := crypto.VerifyHMAC_SHA256(key, mac, b)
	if !valid {
		return keystoreservice.ErrAuthenticationFailed
	}

	r.mux.Lock()

	r.key = key

	r.mux.Unlock()

	return nil
}

func (r *Repository) Update(k *keystore.Keystore) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return fmt.Errorf("authentication required")
	}

	e, err := r.enclaves.Enclave(r.key)
	if err != nil {
		return err
	}

	err = e.UpdateKeystore(k)
	if err != nil {
		return err
	}

	b, err := enclaverepo.MarshalEnclave(e)
	if err != nil {
		return err
	}

	b, err = r.seal(b, r.key, e.Salt())
	if err != nil {
		return err
	}

	return r.enclaves.Write(b)
}

func (r *Repository) startHealthCheck() {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.mux.Lock()
			elapsed := time.Since(r.lastHealthCheck)

			if elapsed < time.Minute {
				r.mux.Unlock()

				continue
			}

			if r.key != nil {
				// TODO: re-enable health check
				//r.key = nil

				//log.Debug("health check failed")
			}

			r.mux.Unlock()
		}
	}
}

func (r *Repository) seal(b, key, salt []byte) ([]byte, error) {
	b, err := crypto.AES256CBC_Encrypt(key, b)
	if err != nil {
		return nil, err
	}

	// authenticate
	mac := crypto.HMAC_SHA256(key, b)

	// prepend salt and mac to the ciphertext
	b = append(salt, append(mac, b...)...)

	return b, nil
}
