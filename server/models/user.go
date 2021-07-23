package models

import (
	"encoding/json"
	"fmt"
	"myst/server/crypto"
	"myst/server/database"
	"myst/server/util"
	"time"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

type User struct {
	ID           string            `json:"id"`
	Username     string            `json:"username"`
	PasswordHash string            `json:"password_hash"`
	Keystores    map[string]string `json:"keystores"`
	CreatedAt    Timestamp         `json:"created_at"`
	UpdatedAt    Timestamp         `json:"updated_at"`
}

type RestUser struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt Timestamp `json:"created_at"`
	UpdatedAt Timestamp `json:"updated_at"`
}

// NewUser creates and saves a new user
func NewUser(username, password string) (*User, error) {
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

// Save saves the user on the database
func (u *User) Save() error {
	now := Timestamp{time.Now()}
	if u.ID == "" {
		u.ID = util.NewUUID()
		u.CreatedAt = now
	}
	if u.Keystores == nil {
		u.Keystores = make(map[string]string)
	}
	u.UpdatedAt = now
	b, err := json.Marshal(u)
	if err != nil {
		return err
	}
	err = database.Save(fmt.Sprintf("data/users/%s.json", u.Username), b)
	if err != nil {
		return err
	}
	return nil
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

// GetUser loads a user from the database and returns it
func GetUser(username string) (*User, error) {
	b, err := database.Load(fmt.Sprintf("data/users/%s.json", username))
	if err == database.ErrNotFound {
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
