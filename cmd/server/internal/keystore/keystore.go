package keystore

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

type Keystore struct {
	ID        string              `bson:"_id"`
	Name      string              `bson:"name"`
	Keystore  []byte              `bson:"keystore"`
	CreatedAt timestamp.Timestamp `bson:"created_at"`
	UpdatedAt timestamp.Timestamp `bson:"updated_at"`
}

// New creates a keystore entry that holds the binary encrypted keystore data
func New(name string, b []byte) (*Keystore, error) {
	store := &Keystore{
		Name:     name,
		Keystore: b,
	}
	err := store.Save()
	if err != nil {
		return nil, err
	}
	return store, nil
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

type RestKeystore struct {
	ID        string              `json:"id"`
	Name      string              `json:"name"`
	Keystore  []byte              `json:"keystore"`
	CreatedAt timestamp.Timestamp `json:"created_at"`
	UpdatedAt timestamp.Timestamp `json:"updated_at"`
}

// ToRest removes sensitive information from the struct
func (store *Keystore) ToRest() *RestKeystore {
	return &RestKeystore{
		ID:        store.ID,
		Name:      store.Name,
		Keystore:  store.Keystore,
		CreatedAt: store.CreatedAt,
		UpdatedAt: store.UpdatedAt,
	}
}

// Get returns a keystore that matches the field/value pairs provided
func Get(field, value string) (*Keystore, error) {
	switch field {
	case "id":
		field = "_id"
	case "name":
		break
	default:
		return nil, ErrInvalidField
	}
	var u *Keystore
	err := database.DB().Collection("keystores").FindOne(context.Background(), bson.M{field: value}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return u, err
}
