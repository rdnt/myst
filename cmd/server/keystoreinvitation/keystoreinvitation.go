package keystoreinvitation

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"myst/database"
	"myst/timestamp"
	"myst/util"
)

var (
	ErrNotFound     = fmt.Errorf("keystore not found")
	ErrInvalidField = fmt.Errorf("invalid field")
)

type KeystoreInvitation struct {
	ID                    string              `bson:"_id"`
	InviterID             string              `bson:"inviter_id"`
	KeystoreID            string              `bson:"keystore_id"`
	InviteeID             string              `bson:"invitee_id"`
	InviterPublicKey      []byte              `bson:"inviter_public_key"`
	InviteePublicKey      []byte              `bson:"invitee_public_key"`
	EncryptedSharedSecret []byte              `bson:"encrypted_shared_secret"`
	CreatedAt             timestamp.Timestamp `bson:"created_at"`
	UpdatedAt             timestamp.Timestamp `bson:"updated_at"`
}

// New creates a keystore entry that holds the binary encrypted keystore data
func New(inviterID, keystoreID, inviteeID string, inviterPublicKey []byte) (*KeystoreInvitation, error) {
	inv := &KeystoreInvitation{
		InviterID:        inviterID,
		KeystoreID:       keystoreID,
		InviteeID:        inviteeID,
		InviterPublicKey: inviterPublicKey,
	}
	err := inv.Save()
	if err != nil {
		return nil, err
	}
	return inv, nil
}

// Save saves the keystore along with the user key on the database
func (inv *KeystoreInvitation) Save() error {
	now := timestamp.New()
	if inv.ID == "" {
		inv.ID = util.NewUUID()
		inv.CreatedAt = now
	}
	inv.UpdatedAt = now

	_, err := database.DB().Collection("keystore_invitations").InsertOne(context.Background(), inv)
	if err != nil {
		return err
	}

	return nil
}

type RestKeystoreInvitation struct {
	ID         string              `json:"id"`
	InviterID  string              `json:"inviter_id"`
	KeystoreID string              `json:"keystore_id"`
	InviteeID  string              `json:"invitee_id"`
	CreatedAt  timestamp.Timestamp `json:"created_at"`
	UpdatedAt  timestamp.Timestamp `json:"updated_at"`
}

// ToRest removes sensitive information from the struct
func (inv *KeystoreInvitation) ToRest() *RestKeystoreInvitation {
	return &RestKeystoreInvitation{
		ID:         inv.ID,
		InviterID:  inv.InviterID,
		KeystoreID: inv.KeystoreID,
		InviteeID:  inv.InviteeID,
		CreatedAt:  inv.CreatedAt,
		UpdatedAt:  inv.UpdatedAt,
	}
}

// Get returns a keystore that matches the field/value pairs provided
func Get(field, value string) (*KeystoreInvitation, error) {
	switch field {
	case "id":
		field = "_id"
	case "TODO@@@":
		break
	default:
		return nil, ErrInvalidField
	}
	var inv *KeystoreInvitation
	err := database.DB().Collection("keystore_invitations").FindOne(context.Background(), bson.M{field: value}).Decode(&inv)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return inv, err
}
