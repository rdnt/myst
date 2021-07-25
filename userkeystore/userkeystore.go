package userkeystore

import (
	"context"
	"fmt"
	"myst/mongo"
	"myst/util"
)

type UserKeystore struct {
	Key      Key
	Keystore Keystore
}

type Key struct {
	ID         string `json:"id" bson:"_id"`
	UserID     string `json:"user_id" bson:"user_id"`
	KeystoreId string `json:"keystore_id" bson:"keystore_id"`
	Key        []byte `json:"key" bson:"key"`
}

type Keystore struct {
	ID       string `json:"id" bson:"_id"`
	Keystore []byte `json:"keystore" bson:"keystore"`
}

// New creates the files of a user's keystore and key from the
// given payload data
func New(userId string, key, keystore []byte) *UserKeystore {
	keystoreId := util.NewUUID()
	return &UserKeystore{
		Key: Key{
			ID:         fmt.Sprintf("%s-%s", userId, keystoreId),
			UserID:     userId,
			KeystoreId: keystoreId,
			Key:        key,
		},
		Keystore: Keystore{
			ID:       keystoreId,
			Keystore: keystore,
		},
	}
}

func (uk *UserKeystore) Save() {
	res, err := mongo.DB().Collection("user_keystores").InsertOne(context.Background(), uk.Key)
	fmt.Println(res, err)

	res, err = mongo.DB().Collection("keystores").InsertOne(context.Background(), uk.Keystore)
	fmt.Println(res, err)

	//fmt.Println("save", litter.Sdump(uk))
}
