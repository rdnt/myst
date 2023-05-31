package remote

import (
	"encoding/json"

	"github.com/pkg/errors"

	"myst/pkg/crypto"
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/user"
	"myst/src/client/enclaverepo"
	"myst/src/server/rest/generated"
)

func InvitationFromJSON(gen generated.Invitation) (invitation.Invitation, error) {
	status, err := invitation.StatusFromString(string(gen.Status))
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "invalid status")
	}

	return invitation.Invitation{
		Id: gen.Id,
		Keystore: keystore.Keystore{
			RemoteId: gen.Keystore.Id,
			Name:     gen.Keystore.Name,
		},
		Inviter: user.User{
			Id:        gen.Inviter.Id,
			Username:  gen.Inviter.Username,
			PublicKey: gen.Inviter.PublicKey,
		},
		Invitee: user.User{
			Id:        gen.Invitee.Id,
			Username:  gen.Invitee.Username,
			PublicKey: gen.Invitee.PublicKey,
		},
		EncryptedKeystoreKey: gen.EncryptedKeystoreKey,
		Status:               status,
		CreatedAt:            gen.CreatedAt,
		UpdatedAt:            gen.CreatedAt,
		AcceptedAt:           gen.AcceptedAt,
		DeclinedAt:           gen.DeclinedAt,
		DeletedAt:            gen.DeletedAt,
	}, nil
}

func KeystoreToJSON(k keystore.Keystore) generated.Keystore {
	return generated.Keystore{
		Id:      k.Id,
		Name:    k.Name,
		Version: k.Version,
	}
}

func UserFromJSON(u generated.User) user.User {
	return user.User{
		Id:        u.Id,
		Username:  u.Username,
		PublicKey: u.PublicKey,
	}
}

func KeystoreFromJSON(gen generated.Keystore, keystoreKey []byte) (keystore.Keystore, error) {
	decryptedKeystorePayload, err := crypto.AES256CBC_Decrypt(keystoreKey, gen.Payload)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "aes256cbc decrypt failed when parsing keystore")
	}

	var jk enclaverepo.KeystoreJSON

	err = json.Unmarshal(decryptedKeystorePayload, &jk)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to unmarshal keystore")
	}

	k, err := enclaverepo.KeystoreFromJSON(jk)
	if err != nil {
		return keystore.Keystore{}, err
	}

	k.ReadOnly = true

	return k, nil
}
