package remote

import (
	"time"

	"github.com/pkg/errors"

	"myst/internal/client/application/domain/invitation"
	"myst/internal/client/application/domain/keystore"
	"myst/internal/server/api/http/generated"
)

func InvitationToJSON(inv invitation.Invitation) generated.Invitation {
	return generated.Invitation{
		Id:          inv.Id,
		InviterId:   inv.InviterId,
		KeystoreId:  inv.KeystoreId,
		InviteeId:   inv.InviteeId,
		InviterKey:  inv.InviterKey,
		InviteeKey:  inv.InviteeKey,
		KeystoreKey: inv.KeystoreKey,
		Status:      inv.Status.String(),
		CreatedAt:   inv.CreatedAt.Unix(),
		UpdatedAt:   inv.UpdatedAt.Unix(),
	}
}

func InvitationFromJSON(gen generated.Invitation) (invitation.Invitation, error) {
	status, err := invitation.StatusFromString(gen.Status)
	if err != nil {
		return invitation.Invitation{}, errors.WithMessage(err, "invalid status")
	}

	return invitation.Invitation{
		Id:          gen.Id,
		InviterId:   gen.InviterId,
		KeystoreId:  gen.KeystoreId,
		InviteeId:   gen.InviteeId,
		InviterKey:  gen.InviterKey,
		InviteeKey:  gen.InviteeKey,
		KeystoreKey: gen.KeystoreKey,
		Status:      status,
		CreatedAt:   time.Unix(gen.CreatedAt, 0),
		UpdatedAt:   time.Unix(gen.CreatedAt, 0),
	}, nil
}

func KeystoreToJSON(k keystore.Keystore) generated.Keystore {

	return generated.Keystore{
		Id:      k.Id,
		Name:    k.Name,
		Version: k.Version,
	}
}

func KeystoreFromJSON(gen generated.Keystore) (keystore.Keystore, error) {
	//CreatedAt
	//Id
	//Name
	//OwnerId
	//Payload
	//UpdatedAt

	// FIXME: this should decode the keystore from the generated one. generated keystore should also return the encrypted
	// payload so that the client can decrypt and parse it into a keystore.Keystore

	return keystore.Keystore{
		Id:      gen.Id,
		Name:    gen.Name,
		Version: gen.Version,
	}, nil
}
