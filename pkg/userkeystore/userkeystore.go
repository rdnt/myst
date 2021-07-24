package userkeystore

import (
	"fmt"
	"github.com/sanity-io/litter"
	"myst/pkg/mongo"
)

type UserKeystore struct {
	UserID string `json:"user_id"`
	Key    []byte `json:"key"`
	Store  []byte `json:"store"`
}

// New creates the files of a user's keystore and key from the
// given payload data
func New(uid string, key, store []byte) *UserKeystore {
	_ = mongo.SaveUserKeystore()
	return &UserKeystore{
		UserID: uid,
		Key:    key,
		Store:  store,
	}
}

func (uk *UserKeystore) Save() {
	fmt.Println("save", litter.Sdump(uk))
}
