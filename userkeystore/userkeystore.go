package userkeystore

import (
	"context"
	"myst/database"
	"myst/timestamp"
	"myst/util"
)

type UserKeystore struct {
	ID         string              `json:"id" bson:"_id"`
	UserID     string              `json:"user-id" bson:"user_id"`
	KeystoreID string              `json:"keystore_id" bson:"keystore_id"`
	Key        []byte              `json:"key" bson:"key"`
	CreatedAt  timestamp.Timestamp `json:"created_at" bson:"created_at"`
	UpdatedAt  timestamp.Timestamp `json:"updated_at" bson:"updated_at"`
}

// New creates the connection of a user to a keystore with the user's key in
// encrypted form as sent by the client
func New(uid, kid string, key []byte) *UserKeystore {
	return &UserKeystore{
		UserID:     uid,
		KeystoreID: kid,
		Key:        key,
	}
}

// Save saves the keystore along with the user key on the database
func (uk *UserKeystore) Save() error {
	now := timestamp.New()
	if uk.ID == "" {
		uk.ID = util.NewUUID()
		uk.CreatedAt = now
	}
	uk.UpdatedAt = now

	_, err := database.DB().Collection("user_keystores").InsertOne(context.Background(), uk)
	if err != nil {
		return err
	}

	return nil
}
