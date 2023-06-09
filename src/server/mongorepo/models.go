package mongorepo

import (
	"time"
)

type User struct {
	Id           string    `bson:"_id"`
	Username     string    `bson:"username"`
	PasswordHash string    `bson:"passwordHash"`
	PublicKey    []byte    `bson:"publicKey"`
	CreatedAt    time.Time `bson:"createdAt"`
	UpdatedAt    time.Time `bson:"updatedAt"`
}

type Keystore struct {
	Id        string    `bson:"_id"`
	Name      string    `bson:"name"`
	Payload   []byte    `bson:"payload"`
	OwnerId   string    `bson:"ownerId"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type Invitation struct {
	Id                   string    `bson:"_id"`
	KeystoreId           string    `bson:"keystoreId"`
	InviterId            string    `bson:"inviterId"`
	InviteeId            string    `bson:"inviteeId"`
	EncryptedKeystoreKey []byte    `bson:"encryptedKeystoreKey"`
	Status               string    `bson:"status"`
	CreatedAt            time.Time `bson:"createdAt"`
	UpdatedAt            time.Time `bson:"updatedAt"`
	AcceptedAt           time.Time `bson:"acceptedAt"`
	DeclinedAt           time.Time `bson:"declinedAt"`
	DeletedAt            time.Time `bson:"deletedAt"`
	CancelledAt          time.Time `bson:"cancelledAt"`
	FinalizedAt          time.Time `bson:"finalizedAt"`
}
