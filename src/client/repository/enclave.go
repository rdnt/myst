package repository

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"

	"github.com/pkg/errors"

	"myst/pkg/crypto"
	"myst/src/client/application/domain/enclave"
	"myst/src/client/application/domain/keystore"
)

func (r *Repository) enclavePath() string {
	return "data.myst"
}

func (r *Repository) enclave(argon2idKey []byte) (*enclave.Enclave, error) {
	f, err := r.fs.Open(r.enclavePath())
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(f)
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
	ks := map[string]Keystore{}

	eks, err := e.Keystores()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get keystores")
	}
	for _, k := range eks {
		ks[k.Id] = KeystoreToJSON(k)
	}

	var jrem *Remote
	rem := e.Remote()

	if rem != nil {
		jrem = &Remote{
			Address:    rem.Address,
			Username:   rem.Username,
			Password:   rem.Password,
			PublicKey:  rem.PublicKey,
			PrivateKey: rem.PrivateKey,
		}
	}

	return json.Marshal(Enclave{
		Keystores: ks,
		Keys:      e.Keys(),
		Remote:    jrem,
	})
}

func enclaveFromJSON(b, salt []byte) (*enclave.Enclave, error) {
	e := &Enclave{}

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

	var rem *enclave.Remote
	jrem := e.Remote

	if jrem != nil {
		rem = &enclave.Remote{
			Address:    jrem.Address,
			Username:   jrem.Username,
			Password:   jrem.Password,
			PublicKey:  jrem.PublicKey,
			PrivateKey: jrem.PrivateKey,
		}
	}

	return enclave.New(
		enclave.WithKeystores(ks),
		enclave.WithSalt(salt),
		enclave.WithRemote(rem),
	)
}

func getSaltFromData(b []byte) ([]byte, error) {
	return b[:crypto.DefaultArgon2IdParams.SaltLength], nil
}
