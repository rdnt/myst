package keystorerepo

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"myst/internal/client/application/domain/enclave"
	"myst/internal/client/application/domain/keystore"
	"myst/pkg/crypto"
)

func (r *Repository) enclavePath() string {
	return path.Join(r.path, "data.myst")
}

func (r *Repository) enclave(argon2idKey []byte) (*enclave.Enclave, error) {
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

	valid := crypto.VerifyHMAC_SHA256(argon2idKey, mac, b)
	if !valid {
		return nil, fmt.Errorf("authentication failed")
	}

	b, err = crypto.AES256CBC_Decrypt(argon2idKey, b)
	if err != nil {
		return nil, err
	}

	return enclaveFromJSON(b, salt)
}

func enclaveToJSON(e *enclave.Enclave) ([]byte, error) {
	ks := map[string]JSONKeystore{}

	for _, k := range e.Keystores() {
		ks[k.Id] = KeystoreToJSON(k)
	}

	return json.Marshal(JSONEnclave{Keystores: ks, Keys: e.Keys()})
}

func enclaveFromJSON(b, salt []byte) (*enclave.Enclave, error) {
	e := &JSONEnclave{}

	err := json.Unmarshal(b, e)
	if err != nil {
		return nil, err
	}

	ks := map[string]keystore.Keystore{}

	for _, k := range e.Keystores {
		k, err := KeystoreFromJSON(k)
		if err != nil {
			return nil, err
		}

		ks[k.Id] = k
	}

	return enclave.New(
		enclave.WithKeystores(ks),
		enclave.WithSalt(salt),
	)
}

func getSaltFromData(b []byte) ([]byte, error) {
	return b[:crypto.DefaultArgon2IdParams.SaltLength], nil
}
