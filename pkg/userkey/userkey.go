package userkey

import (
	"context"
	"fmt"

	"myst/pkg/mongo"
	"myst/pkg/util"

	"go.mongodb.org/mongo-driver/bson"

	"myst/pkg/timestamp"
)

var (
	ErrNotFound     = fmt.Errorf("user_key not found")
	ErrInvalidField = fmt.Errorf("invalid field")
)

type UserKey struct {
	ID         string              `bson:"_id"`
	UserID     string              `bson:"user_id"`
	KeystoreID string              `bson:"keystore_id"`
	Key        []byte              `bson:"key"`
	CreatedAt  timestamp.Timestamp `bson:"created_at"`
	UpdatedAt  timestamp.Timestamp `bson:"updated_at"`
}

// New creates the connection of a user to a domain with the user's key in
// encrypted form as sent by the client
func New(uid, kid string, key []byte) (*UserKey, error) {
	uk := &UserKey{
		UserID:     uid,
		KeystoreID: kid,
		Key:        key,
	}
	err := uk.Save()
	if err != nil {
		return nil, err
	}
	return uk, nil
}

// Save saves the domain along with the user key on the database
func (uk *UserKey) Save() error {
	now := timestamp.New()
	if uk.ID == "" {
		uk.ID = util.NewUUID()
		uk.CreatedAt = now
	}
	uk.UpdatedAt = now

	_, err := mongo.DB().Collection("user_keys").InsertOne(context.Background(), uk)
	if err != nil {
		return err
	}

	return nil
}

type RestUserKey struct {
	ID         string              `json:"id"`
	UserID     string              `json:"user_id"`
	KeystoreID string              `json:"keystore_id"`
	Key        []byte              `json:"key"`
	CreatedAt  timestamp.Timestamp `json:"created_at"`
	UpdatedAt  timestamp.Timestamp `json:"updated_at"`
}

// ToRest removes sensitive information from the struct
func (uk *UserKey) ToRest() *RestUserKey {
	return &RestUserKey{
		ID:         uk.ID,
		UserID:     uk.UserID,
		KeystoreID: uk.KeystoreID,
		Key:        uk.Key,
		CreatedAt:  uk.CreatedAt,
		UpdatedAt:  uk.UpdatedAt,
	}
}

// Get returns a domain that matches the field/value pairs provided
func Get(filters map[string]string) (*UserKey, error) {
	var u *UserKey

	m := bson.M{}
	for k, v := range filters {
		switch k {
		case "id":
			k = "_id"
		case "user_id", "keystore_id":
			break
		default:
			return nil, ErrInvalidField
		}
		m[k] = v
	}

	err := mongo.DB().Collection("user_keys").FindOne(context.Background(), m).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return u, err
}
