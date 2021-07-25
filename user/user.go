package user

import (
	"context"
	"encoding/json"
	"fmt"
	"myst/crypto"
	"myst/mongo"
	"myst/storage"
	"myst/timestamp"
	"myst/util"
)

var (
	ErrNotFound = fmt.Errorf("user not found")
)

type User struct {
	ID           string              `json:"id"`
	Username     string              `json:"username"`
	PasswordHash string              `json:"password_hash"`
	Keystores    map[string]string   `json:"keystores"`
	CreatedAt    timestamp.Timestamp `json:"created_at"`
	UpdatedAt    timestamp.Timestamp `json:"updated_at"`
}

// Save saves the user on the storage
func (u *User) Save() error {
	now := timestamp.New()
	if u.ID == "" {
		u.ID = util.NewUUID()
		u.CreatedAt = now
	}
	if u.Keystores == nil {
		u.Keystores = make(map[string]string)
	}
	u.UpdatedAt = now

	_, err := mongo.DB().Collection("users").InsertOne(context.Background(), u)
	if err != nil {
		return err
	}
	return nil
}

type RestUser struct {
	ID        string              `json:"id"`
	Username  string              `json:"username"`
	CreatedAt timestamp.Timestamp `json:"created_at"`
	UpdatedAt timestamp.Timestamp `json:"updated_at"`
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
	hash, err := crypto.HashPassword(password)
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
func Get(username string) (*User, error) {
	b, err := storage.Load(fmt.Sprintf("data/users/%s.json", username))
	if err == storage.ErrNotFound {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	var u *User
	err = json.Unmarshal(b, &u)
	if err != nil {
		return nil, err
	}

	return u, nil
}
