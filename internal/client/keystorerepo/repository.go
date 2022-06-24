package keystorerepo

import (
	"crypto/sha256"
	"fmt"
	"os"
	"sync"
	"time"

	"myst/internal/client/application/domain/enclave"
	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/application/keystoreservice"
	"myst/pkg/crypto"
	"myst/pkg/logger"

	"github.com/pkg/errors"
)

var log = logger.New("keystorerepo", logger.Green)

type Repository struct {
	mux             sync.Mutex
	path            string
	key             []byte
	lastHealthCheck time.Time
	//remote          remote.Remote
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

func (r *Repository) Keystore(id string) (keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return keystore.Keystore{}, fmt.Errorf("authentication required")
	}

	e, err := r.enclave(r.key)
	if err != nil {
		return keystore.Keystore{}, err
	}

	k, err := e.Keystore(id)
	if err != nil {
		return keystore.Keystore{}, err
	}

	return k, nil
}

//func (r *Repository) Sync() error {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//
//	if r.key == nil {
//		return fmt.Errorf("authentication required")
//	}
//
//	rks, err := r.remote.Keystores()
//	if err != nil {
//		return errors.WithMessage(err, "failed to get remote keystores")
//	}
//
//	ks, err := r.Keystores()
//	if err != nil {
//		return errors.WithMessage(err, "failed to get local keystores")
//	}
//
//	for _, k := range rks {
//		if _, ok := ks[k.Id]; !ok {
//			k, err = r.createKeystore(k)
//			if err != nil {
//				return errors.WithMessage(err, "failed to create keystore")
//			}
//
//			ks[k.Id] = k
//		} else {
//			err = r.updateKeystore(k)
//			if err != nil {
//				return errors.WithMessage(err, "failed to update keystore")
//			}
//		}
//	}
//
//	return nil
//}

func (r *Repository) Keystores() (map[string]keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, err := os.ReadFile(r.enclavePath())
	if errors.Is(err, os.ErrNotExist) {
		return nil, keystoreservice.ErrInitializationRequired
	} else if err != nil {
		return nil, err
	}

	if r.key == nil {
		return nil, keystoreservice.ErrAuthenticationRequired
	}

	e, err := r.enclave(r.key)
	if err != nil {
		return nil, err
	}

	return e.Keystores()
}

func (r *Repository) DeleteKeystore(id string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return fmt.Errorf("authentication required")
	}

	e, err := r.enclave(r.key)
	if err != nil {
		return err
	}

	err = e.DeleteKeystore(id)
	if err != nil {
		return err
	}

	err = r.sealAndWrite(e)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CreateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return keystore.Keystore{}, fmt.Errorf("authentication required")
	}

	return r.createKeystore(k)
}

func (r *Repository) createKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	e, err := r.enclave(r.key)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "failed to get enclave")
	}

	err = e.AddKeystore(k)
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

	p := crypto.DefaultArgon2IdParams

	salt, err := crypto.GenerateRandomBytes(uint(p.SaltLength))
	if err != nil {
		return err
	}

	// TODO: check if enclave is already created, if it is, return error
	e, err := enclave.New(enclave.WithSalt(salt))
	if err != nil {
		return err
	}

	key := crypto.Argon2Id([]byte(password), salt)
	if err != nil {
		return err
	}

	r.key = key

	err = r.sealAndWrite(e)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateSyncKeypair(publicKey, privateKey []byte) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return fmt.Errorf("authentication required")
	}

	e, err := r.enclave(r.key)
	if err != nil {
		return err
	}

	e.SetSyncKeypair(publicKey, privateKey)

	return r.sealAndWrite(e)
}

func (r *Repository) Authenticate(password string) error {
	r.mux.Lock()

	b, err := os.ReadFile(r.enclavePath())
	if err != nil {
		r.mux.Unlock()
		return err
	}

	r.mux.Unlock()

	salt, err := getSaltFromData(b)
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

func (r *Repository) UpdateKeystore(k keystore.Keystore) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return fmt.Errorf("authentication required")
	}

	return r.updateKeystore(k)
}

func (r *Repository) updateKeystore(k keystore.Keystore) error {
	e, err := r.enclave(r.key)
	if err != nil {
		return err
	}

	err = e.UpdateKeystore(k)
	if err != nil {
		return err
	}

	return r.sealAndWrite(e)
}

//func (r *Repository) KeystoreKey(keystoreId string) ([]byte, error) {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//
//	if r.key == nil {
//		return nil, fmt.Errorf("authentication required")
//	}
//
//	e, err := r.enclave(r.key)
//	if err != nil {
//		return nil, err
//	}
//
//	keys := e.Keys()
//	b, ok := keys[keystoreId]
//	if !ok {
//		return nil, fmt.Errorf("keystore key not found")
//	}
//
//	return b, nil
//}

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

func (r *Repository) sealAndWrite(e *enclave.Enclave) error {
	b, err := enclaveToJSON(e)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal enclave")
	}

	b, err = crypto.AES256CBC_Encrypt(r.key, b)
	if err != nil {
		return err
	}

	// authenticate
	mac := crypto.HMAC_SHA256(r.key, b)

	// prepend salt and mac to the ciphertext
	b = append(e.Salt(), append(mac, b...)...)

	return os.WriteFile(r.enclavePath(), b, 0600)
}

func (r *Repository) SyncKeypair() (publicKey, privateKey []byte, err error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return nil, nil, fmt.Errorf("authentication required")
	}

	e, err := r.enclave(r.key)
	if err != nil {
		return nil, nil, err
	}

	return e.SyncKeypair()
}
