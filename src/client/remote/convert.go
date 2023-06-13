package remote

import (
	"encoding/json"
	"myst/src/client/application/domain/entry"
	"time"

	"github.com/pkg/errors"

	"myst/pkg/crypto"
	"myst/src/client/application/domain/invitation"
	"myst/src/client/application/domain/keystore"
	"myst/src/client/application/domain/user"
	"myst/src/server/rest/generated"
)

type KeystoreJSON struct {
	Id        string      `json:"id"`
	RemoteId  string      `json:"remoteId"`
	Name      string      `json:"name"`
	Version   int         `json:"version"`
	Entries   []EntryJSON `json:"entries"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type EntryJSON struct {
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

//func KeystoreToJSON(k keystore.keystore) generated.keystore {
//	return generated.keystore{
//		Id:      k.Id,
//		Name:    k.Name,
//		Version: k.Version,
//	}
//}
//
//func UserFromJSON(u generated.User) user.User {
//	return user.User{
//		Id:        u.Id,
//		Username:  u.Username,
//		PublicKey: u.PublicKey,
//	}
//}

func KeystoreFromJSON(gen generated.Keystore, keystoreKey []byte) (keystore.Keystore, error) {
	decryptedKeystorePayload, err := crypto.AES256CBC_Decrypt(keystoreKey, gen.Payload)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "aes256cbc decrypt failed when parsing keystore")
	}

	var jk KeystoreJSON

	err = json.Unmarshal(decryptedKeystorePayload, &jk)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to unmarshal keystore")
	}

	k, err := keystoreFromJSON(jk)
	if err != nil {
		return keystore.Keystore{}, err
	}

	k.ReadOnly = true
	k.Key = keystoreKey

	return k, nil
}

func keystoreToJSON(k keystore.Keystore) KeystoreJSON {
	entries := []EntryJSON{}

	for _, e := range k.Entries {
		entries = append(entries, EntryJSON{
			Id:        e.Id,
			Website:   e.Website,
			Username:  e.Username,
			Password:  e.Password,
			Notes:     e.Notes,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}

	return KeystoreJSON{
		Id:        k.Id,
		RemoteId:  k.RemoteId,
		Name:      k.Name,
		Version:   k.Version,
		Entries:   entries,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func keystoreFromJSON(k KeystoreJSON) (keystore.Keystore, error) {
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
	}, nil
}
