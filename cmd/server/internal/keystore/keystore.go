package keystore

import (
	"context"
	"myst/database"
	"myst/timestamp"
	"myst/util"
)

type Keystore struct {
	ID        string              `json:"id" bson:"_id"`
	Keystore  []byte              `json:"keystore" bson:"keystore"`
	CreatedAt timestamp.Timestamp `json:"created_at" bson:"created_at"`
	UpdatedAt timestamp.Timestamp `json:"updated_at" bson:"updated_at"`
}

// New creates a keystore entry that holds the binary encrypted keystore data
func New(store []byte) *Keystore {
	return &Keystore{
		ID:       "",
		Keystore: store,
	}
}

// Save saves the keystore along with the user key on the database
func (store *Keystore) Save() error {
	now := timestamp.New()
	if store.ID == "" {
		store.ID = util.NewUUID()
		store.CreatedAt = now
	}
	store.UpdatedAt = now

	_, err := database.DB().Collection("keystores").InsertOne(context.Background(), store)
	if err != nil {
		return err
	}

	return nil
}
