package enclaverepo

import (
	"crypto/sha256"
	"fmt"
	enclave2 "myst/src/client/enclaverepo/enclave"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"

	"myst/pkg/crypto"
	"myst/pkg/logger"
	"myst/src/client/application"
	"myst/src/client/application/domain/credentials"
	"myst/src/client/application/domain/keystore"
)

var log = logger.New("repository", logger.Green)

type Repository struct {
	mux             sync.Mutex
	key             []byte
	lastHealthCheck time.Time
	path            string
	// remote          remote.Credentials
}

func New(path string) (*Repository, error) {
	r := &Repository{
		path: path,
	}

	go r.startHealthCheck()

	// TODO: remove simulated health check
	// go func() {
	//	for {
	//		time.Sleep(10 * time.Second)
	//		r.HealthCheck()
	//	}
	// }()

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

// func (r *Repository) Sync() error {
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
// }

func (r *Repository) Keystores() (map[string]keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, err := os.ReadFile(r.enclavePath())
	if errors.Is(err, os.ErrNotExist) {
		return nil, application.ErrInitializationRequired
	} else if err != nil {
		return nil, err
	}

	if r.key == nil {
		return nil, application.ErrAuthenticationRequired
	}

	e, err := r.enclave(r.key)
	if err != nil {
		return nil, err
	}

	return e.Keystores()
}

func (r *Repository) IsInitialized() error {
	r.mux.Lock()
	defer r.mux.Unlock()

	_, err := os.ReadFile(r.enclavePath())
	if errors.Is(err, os.ErrNotExist) {
		return application.ErrInitializationRequired
	} else if err != nil {
		return err
	}

	if r.key == nil {
		return application.ErrAuthenticationRequired
	}

	return nil
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

// CreateEnclave initializes the enclave with the given password
func (r *Repository) Initialize(password string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	p := crypto.DefaultArgon2IdParams

	salt, err := crypto.GenerateRandomBytes(uint(p.SaltLength))
	if err != nil {
		return err
	}

	key := crypto.Argon2Id([]byte(password), salt)
	if err != nil {
		return err
	}

	// TODO: check if enclave is already created, if it is, return error
	e, err := enclave2.New(
		enclave2.WithSalt(salt),
	)
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
		return application.ErrAuthenticationFailed
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

// func (r *Repository) EncryptedKeystoreKey(keystoreId string) ([]byte, error) {
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
// }

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
				// r.key = nil

				// log.Debug("health check failed")
			}

			r.mux.Unlock()
		}
	}
}

func (r *Repository) sealAndWrite(e *enclave2.Enclave) error {
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

// func (r *Repository) Keypair() (publicKey, privateKey []byte, err error) {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//
//	if r.key == nil {
//		return nil, nil, fmt.Errorf("authentication required")
//	}
//
//	e, err := r.enclave(r.key)
//	if err != nil {
//		return nil, nil, err
//	}
//
//	publicKey, privateKey, err = e.Keypair()
//	if err == enclave.ErrNotSet {
//		return nil, nil, ErrNotSet
//	} else if err != nil {
//		return nil, nil, err
//	}
//
//	return publicKey, privateKey, nil
// }
//
// func (r *Repository) SetKeypair(publicKey, privateKey []byte) error {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//
//	if r.key == nil {
//		return fmt.Errorf("authentication required")
//	}
//
//	e, err := r.enclave(r.key)
//	if err != nil {
//		return err
//	}
//
//	e.SetKeypair(publicKey, privateKey)
//
//	return r.sealAndWrite(e)
// }

// func (r *Repository) UserInfo() (username, password string, err error) {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//
//	if r.key == nil {
//		return "", "", fmt.Errorf("authentication required")
//	}
//
//	e, err := r.enclave(r.key)
//	if err != nil {
//		return "", "", err
//	}
//
//	username, password, err = e.UserInfo()
//	if err == enclave.ErrNotSet {
//		return "", "", ErrNotSet
//	} else if err != nil {
//		return "", "", err
//	}
//
//	return username, password, nil
// }
//
// func (r *Repository) SetUserInfo(username, password string) error {
//	r.mux.Lock()
//	defer r.mux.Unlock()
//
//	if r.key == nil {
//		return fmt.Errorf("authentication required")
//	}
//
//	e, err := r.enclave(r.key)
//	if err != nil {
//		return err
//	}
//
//	e.SetUserInfo(username, password)
//
//	return r.sealAndWrite(e)
// }

func (r *Repository) SetCredentials(address, username, password string, publicKey, privateKey []byte) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return fmt.Errorf("authentication required")
	}

	e, err := r.enclave(r.key)
	if err != nil {
		return err
	}

	e.SetRemote(credentials.Credentials{
		Address:    address,
		Username:   username,
		Password:   password,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	})

	return r.sealAndWrite(e)
}

func (r *Repository) Credentials() (credentials.Credentials, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return credentials.Credentials{}, fmt.Errorf("authentication required")
	}

	e, err := r.enclave(r.key)
	if err != nil {
		return credentials.Credentials{}, err
	}

	rem := e.Remote()
	if rem == nil {
		return credentials.Credentials{}, enclave2.ErrRemoteNotSet
	}

	return *rem, nil
}
