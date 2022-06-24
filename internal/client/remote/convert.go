package remote

import (
	"encoding/json"

	"github.com/pkg/errors"

	"myst/internal/client/application/domain/invitation"
	"myst/internal/client/application/domain/keystore"
	"myst/internal/client/keystorerepo"
	"myst/internal/server/api/http/generated"
	"myst/pkg/crypto"
	"myst/pkg/logger"
)

func InvitationFromJSON(gen generated.Invitation) (invitation.Invitation, error) {
	status, err := invitation.StatusFromString(string(gen.Status))
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "invalid status")
	}

	return invitation.Invitation{
		Id:           gen.Id,
		InviterId:    gen.InviterId,
		KeystoreId:   gen.KeystoreId,
		KeystoreName: gen.KeystoreName,
		InviteeId:    gen.InviteeId,
		InviterKey:   gen.InviterKey,
		InviteeKey:   gen.InviteeKey,
		KeystoreKey:  gen.KeystoreKey,
		Status:       status,
		CreatedAt:    gen.CreatedAt,
		UpdatedAt:    gen.CreatedAt,
		AcceptedAt:   gen.AcceptedAt,
		DeclinedAt:   gen.DeclinedAt,
		DeletedAt:    gen.DeletedAt,
	}, nil
}

func KeystoreToJSON(k keystore.Keystore) generated.Keystore {

	return generated.Keystore{
		Id:      k.Id,
		Name:    k.Name,
		Version: k.Version,
	}
}

func KeystoreFromJSON(gen generated.Keystore, keystoreKey []byte) (keystore.Keystore, error) {
	//CreatedAt
	//Id
	//Name
	//OwnerId
	//Payload
	//UpdatedAt

	logger.Error("@@@@@@@@@@@@@@@@@@@@@@@@@@", string(keystoreKey))

	decryptedKeystorePayload, err := crypto.AES256CBC_Decrypt(keystoreKey, gen.Payload)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "aes256cbc decrypt failed when parsing keystore")
	}

	logger.Error("DECRYPTEEEED", string(decryptedKeystorePayload))

	var jk keystorerepo.Keystore

	err = json.Unmarshal(decryptedKeystorePayload, &jk)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to unmarshal keystore")
	}

	return keystorerepo.KeystoreFromJSON(jk)

	// FIXME: this should decode the keystore from the generated one. generated keystore should also return the encrypted
	// payload so that the client can decrypt and parse it into a keystore.Keystore

}
