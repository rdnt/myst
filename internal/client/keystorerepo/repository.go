package keystorerepo

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"myst/internal/client/application/keystoreservice"
	"myst/pkg/logger"
	"os"
	"path"
	"sync"
	"time"

	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/keystorerepo/enclave"
	"myst/pkg/crypto"
)

var log = logger.New("keystorerepo", logger.Green)

type Enclave struct {
	Keystores map[string]JSONKeystore `json:"keystores"`
	Keys      map[string][]byte       `json:"keys"`
}

type Repository struct {
	mux             sync.Mutex
	path            string
	key             []byte
	lastHealthCheck time.Time
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

func (r *Repository) Keystore(id string) (*keystore.Keystore, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return nil, fmt.Errorf("authentication required")
	}

	e, err := r.enclave()
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

	_, err := os.ReadFile(r.enclavePath())
	if errors.Is(err, os.ErrNotExist) {
		return nil, keystoreservice.ErrInitializationRequired
	} else if err != nil {
		return nil, err
	}

	if r.key == nil {
		return nil, keystoreservice.ErrAuthenticationRequired
	}

	e, err := r.enclave()
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

	e, err := r.enclave()
	if err != nil {
		return err
	}

	err = e.DeleteKeystore(id)
	if err != nil {
		return err
	}

	b, err := marshalEnclave(e)
	if err != nil {
		return err
	}

	err = r.sealAndWrite(b, r.key, e.Salt())
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

	e, err := r.enclave()
	if err != nil {
		return nil, err
	}

	err = e.AddKeystore(k)
	if err != nil {
		return nil, err
	}

	b, err := marshalEnclave(e)
	if err != nil {
		return nil, err
	}

	err = r.sealAndWrite(b, r.key, e.Salt())
	if err != nil {
		return nil, err
	}

	return k, nil
}

func (r *Repository) Initialize(password string) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	// TODO: check if enclave is already created, if it is, return error
	e, err := enclave.New()
	if err != nil {
		return err
	}

	b, err := marshalEnclave(e)
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

	err = r.sealAndWrite(b, key, salt)
	if err != nil {
		return err
	}

	r.key = key

	return nil
}

func (r *Repository) Authenticate(password string) error {
	r.mux.Lock()

	b, err := os.ReadFile(r.enclavePath())
	if err != nil {
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

func (r *Repository) enclave() (*enclave.Enclave, error) {
	b, err := os.ReadFile(r.enclavePath())
	if err != nil {
		return nil, err
	}

	salt, err := getSaltFromData(b)
	if err != nil {
		return nil, err
	}

	p := crypto.DefaultArgon2IdParams

	mac := b[p.SaltLength : sha256.Size+p.SaltLength]
	b = b[p.SaltLength+sha256.Size:]

	valid := crypto.VerifyHMAC_SHA256(r.key, mac, b)
	if !valid {
		return nil, fmt.Errorf("authentication failed")
	}

	b, err = crypto.AES256CBC_Decrypt(r.key, b)
	if err != nil {
		return nil, err
	}

	return unmarshalEnclave(b, salt)
}

func (r *Repository) Update(k *keystore.Keystore) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.key == nil {
		return fmt.Errorf("authentication required")
	}

	e, err := r.enclave()
	if err != nil {
		return err
	}

	err = e.UpdateKeystore(k)
	if err != nil {
		return err
	}

	b, err := marshalEnclave(e)
	if err != nil {
		return err
	}

	err = r.sealAndWrite(b, r.key, e.Salt())
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) startHealthCheck() {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("health check!")

			r.mux.Lock()
			elapsed := time.Since(r.lastHealthCheck)

			if elapsed < time.Minute {
				r.mux.Unlock()

				continue
			}

			if r.key != nil {
				r.key = nil

				log.Debug("health check failed")
			}

			r.mux.Unlock()
		}
	}
}

func (r *Repository) sealAndWrite(b, key, salt []byte) error {
	b, err := crypto.AES256CBC_Encrypt(key, b)
	if err != nil {
		return err
	}

	// authenticate
	mac := crypto.HMAC_SHA256(key, b)

	// prepend salt and mac to the ciphertext
	b = append(salt, append(mac, b...)...)

	err = os.WriteFile(r.enclavePath(), b, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) enclavePath() string {
	return path.Join(r.path, "data.myst")
}

func marshalEnclave(e *enclave.Enclave) ([]byte, error) {
	ks := map[string]JSONKeystore{}

	for _, k := range e.Keystores() {
		ks[k.Id()] = ToJSONKeystore(k)
	}

	return json.Marshal(Enclave{Keystores: ks, Keys: e.Keys()})
}

func unmarshalEnclave(b, salt []byte) (*enclave.Enclave, error) {
	e := &Enclave{}

	err := json.Unmarshal(b, e)
	if err != nil {
		return nil, err
	}

	ks := map[string]*keystore.Keystore{}

	for _, k := range e.Keystores {
		k, err := ToKeystore(k)
		if err != nil {
			return nil, err
		}

		ks[k.Id()] = k
	}

	return enclave.New(
		enclave.WithKeystores(ks),
		enclave.WithSalt(salt),
	)
}

func getSaltFromData(b []byte) ([]byte, error) {
	return b[:crypto.DefaultArgon2IdParams.SaltLength], nil
}
