package user

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	crypto2 "myst/pkg/crypto"
	database2 "myst/pkg/database"
	timestamp2 "myst/pkg/timestamp"
	util2 "myst/pkg/util"
)

var (
	ErrNotFound     = fmt.Errorf("user not found")
	ErrInvalidField = fmt.Errorf("invalid field")
)

type User struct {
	ID           string               `bson:"_id"`
	Username     string               `bson:"username"`
	PasswordHash string               `bson:"password_hash"`
	CreatedAt    timestamp2.Timestamp `bson:"created_at"`
	UpdatedAt    timestamp2.Timestamp `bson:"updated_at"`
}

// Save saves the user on the storage
func (u *User) Save() error {
	now := timestamp2.New()
	if u.ID == "" {
		u.ID = util2.NewUUID()
		u.CreatedAt = now
	}
	u.UpdatedAt = now

	_, err := database2.DB().Collection("users").InsertOne(context.Background(), u)
	if err != nil {
		return err
	}
	return nil
}

type RestUser struct {
	ID        string               `json:"id"`
	Username  string               `json:"username"`
	CreatedAt timestamp2.Timestamp `json:"created_at"`
	UpdatedAt timestamp2.Timestamp `json:"updated_at"`
}

// ToRest removes sensitive information from the struct
func (u *User) ToRest() *RestUser {
	return &RestUser{
		ID:        u.ID,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// New creates and saves a new user
func New(username, password string) (*User, error) {
	hash, err := crypto2.HashPassword(password)
	if err != nil {
		return nil, err
	}
	u := &User{
		Username:     username,
		PasswordHash: hash,
	}
	err = u.Save()
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Get loads a user from the storage and returns it
//func Get(id string) (*User, error) {
//	u := &User{
//		ID: id,
//	}
//	err := database.DB().Collection("users").FindOne(context.Background(), &u).Err()
//	if err == mongo.ErrNoDocuments {
//		return nil, ErrNotFound
//	} else if err != nil {
//		return nil, err
//	}
//	return u, err
//}

// Get returns a user that matches the field/value pairs provided
func Get(field, value string) (*User, error) {
	switch field {
	case "id":
		field = "_id"
	case "username":
		break
	default:
		return nil, ErrInvalidField
	}
	var u *User
	err := database2.DB().Collection("users").FindOne(context.Background(), bson.M{field: value}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return u, err
}
