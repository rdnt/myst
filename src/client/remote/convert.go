package remote

import (
	"encoding/json"

	"github.com/pkg/errors"

	"myst/pkg/crypto"
	"myst/pkg/logger"
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/user"
	"myst/src/client/repository"
	"myst/src/server/rest/generated"
)

func InvitationFromJSON(gen generated.Invitation) (invitation.Invitation, error) {
	status, err := invitation.StatusFromString(string(gen.Status))
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "invalid status")
	}

	return invitation.Invitation{
		Id:                   gen.Id,
		InviterId:            gen.InviterId,
		KeystoreId:           gen.KeystoreId,
		KeystoreName:         gen.KeystoreName,
		InviteeId:            gen.InviteeId,
		InviterPublicKey:     gen.InviterPublicKey,
		InviteePublicKey:     gen.InviteePublicKey,
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
		Id:       u.Id,
		Username: u.Username,
	}
}

func KeystoreFromJSON(gen generated.Keystore, keystoreKey []byte) (keystore.Keystore, error) {
	// CreatedAt
	// Id
	// Name
	// OwnerId
	// Payload
	// UpdatedAt

	logger.Error("@@@@@@@@@@@@@@@@@@@@@@@@@@", string(keystoreKey))

	decryptedKeystorePayload, err := crypto.AES256CBC_Decrypt(keystoreKey, gen.Payload)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "aes256cbc decrypt failed when parsing keystore")
	}

	logger.Error("DECRYPTEEEED", string(decryptedKeystorePayload))

	var jk repository.Keystore

	err = json.Unmarshal(decryptedKeystorePayload, &jk)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to unmarshal keystore")
	}

	return repository.KeystoreFromJSON(jk)

	// FIXME: this should decode the keystore from the generated one. generated keystore should also return the encrypted
	// payload so that the client can decrypt and parse it into a keystore.Keystore

}
