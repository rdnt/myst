package remote

import (
	"encoding/json"
	"time"

	"github.com/samber/lo"

	"myst/src/client/application/domain/entry"
	"myst/src/client/application/domain/user"

	"github.com/pkg/errors"

	"myst/pkg/crypto"
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/server/rest/generated"
)

type keystoreJSON struct {
	Id        string      `json:"id"`
	RemoteId  string      `json:"remoteId"`
	Name      string      `json:"name"`
	Version   int         `json:"version"`
	Entries   []entryJSON `json:"entries"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type entryJSON struct {
	Id        string    `json:"id"`
	Website   string    `json:"website"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func InvitationFromJSON(gen generated.Invitation) (invitation.Invitation, error) {
	status, err := invitation.StatusFromString(string(gen.Status))
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "invalid invitation status")
	}

	return invitation.Invitation{
		Id:               gen.Id,
		RemoteKeystoreId: gen.Keystore.Id,
		KeystoreName:     gen.Keystore.Name,
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
		DeletedAt:            gen.DeletedAt,
		DeclinedAt:           gen.DeclinedAt,
		AcceptedAt:           gen.AcceptedAt,
		CancelledAt:          gen.CancelledAt,
		FinalizedAt:          gen.FinalizedAt,
	}, nil
}

func KeystoreFromJSON(gen generated.Keystore, keystoreKey []byte) (keystore.Keystore, error) {
	decryptedKeystorePayload, err := crypto.AES256CBC_Decrypt(keystoreKey, gen.Payload)
	if err != nil {
		return keystore.Keystore{}, errors.WithMessage(err, "aes256cbc decrypt failed")
	}

	var jk keystoreJSON

	err = json.Unmarshal(decryptedKeystorePayload, &jk)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to unmarshal keystore")
	}

	k := keystoreFromJSON(jk)

	k.ReadOnly = true
	k.Key = keystoreKey

	return k, nil
}

func keystoreToJSON(k keystore.Keystore) keystoreJSON {
	entries := make([]entryJSON, len(k.Entries))

	for i, e := range lo.Values(k.Entries) {
		entries[i] = entryJSON{
			Id:        e.Id,
			Website:   e.Website,
			Username:  e.Username,
			Password:  e.Password,
			Notes:     e.Notes,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		}
	}

	return keystoreJSON{
		Id:        k.Id,
		RemoteId:  k.RemoteId,
		Name:      k.Name,
		Version:   k.Version,
		Entries:   entries,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func keystoreFromJSON(k keystoreJSON) keystore.Keystore {
	entries := make(map[string]entry.Entry, len(k.Entries))

	for _, e := range k.Entries {
		e := entry.Entry{
			Id:        e.Id,
			Website:   e.Website,
			Username:  e.Username,
			Password:  e.Password,
			Notes:     e.Notes,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		}

		entries[e.Id] = e
	}

	return keystore.Keystore{
		Id:        k.Id,
		RemoteId:  k.RemoteId,
		Name:      k.Name,
		Version:   k.Version,
		Entries:   entries,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}
