package mongorepo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"myst/src/server/application/domain/keystore"
)

type Keystore struct {
	Id        string    `bson:"_id"`
	Name      string    `bson:"name"`
	Payload   []byte    `bson:"payload"`
	OwnerId   string    `bson:"ownerId"`
	ViewerIds []string  `bson:"viewerIds"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

func KeystoreToBSON(k keystore.Keystore) Keystore {
	return Keystore{
		Id:        k.Id,
		Name:      k.Name,
		Payload:   k.Payload,
		OwnerId:   k.OwnerId,
		ViewerIds: k.ViewerIds,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func KeystoreFromBSON(k Keystore) keystore.Keystore {
	return keystore.Keystore{
		Id:        k.Id,
		Name:      k.Name,
		Payload:   k.Payload,
		OwnerId:   k.OwnerId,
		ViewerIds: k.ViewerIds,
		CreatedAt: k.CreatedAt,
		UpdatedAt: k.UpdatedAt,
	}
}

func (r *Repository) CreateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	collection := r.mdb.Database(r.database).Collection("keystores")

	kb := KeystoreToBSON(k)

	_, err := collection.InsertOne(context.Background(), kb)
	if err != nil {
		return keystore.Keystore{}, err
	}

	k = KeystoreFromBSON(kb)

	return k, nil
}

func (r *Repository) Keystore(id string) (keystore.Keystore, error) {
	collection := r.mdb.Database(r.database).Collection("keystores")

	res := collection.FindOne(context.Background(), bson.D{{"_id", id}})
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return keystore.Keystore{}, keystore.ErrNotFound
	} else if err != nil {
		return keystore.Keystore{}, err
	}

	var k Keystore
	err = res.Decode(&k)
	if err != nil {
		return keystore.Keystore{}, err
	}

	return KeystoreFromBSON(k), nil
}

func (r *Repository) Keystores() ([]keystore.Keystore, error) {
	collection := r.mdb.Database(r.database).Collection("keystores")

	ctx := context.Background()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var k Keystore
	keystores := []keystore.Keystore{}
	for cur.Next(ctx) {
		err := cur.Decode(&k)
		if err != nil {
			return nil, err
		}

		keystores = append(keystores, KeystoreFromBSON(k))
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return keystores, nil
}

func (r *Repository) UpdateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	collection := r.mdb.Database(r.database).Collection("keystores")

	ctx := context.Background()

	kb := KeystoreToBSON(k)
	res := collection.FindOneAndReplace(ctx, bson.D{{"_id", kb.Id}}, kb)
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return keystore.Keystore{}, keystore.ErrNotFound
	} else if err != nil {
		return keystore.Keystore{}, err
	}

	return KeystoreFromBSON(kb), nil
}

func (r *Repository) DeleteKeystore(id string) error {
	collection := r.mdb.Database(r.database).Collection("keystores")
	ctx := context.Background()

	res := collection.FindOneAndDelete(ctx, bson.D{{"_id", id}})
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return keystore.ErrNotFound
	} else if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UserKeystores(userId string) ([]keystore.Keystore, error) {
	userKeystores := []keystore.Keystore{}

	allKeystores, err := r.Keystores()
	if err != nil {
		return nil, err
	}

	for _, k := range allKeystores {
		if k.OwnerId == userId {
			userKeystores = append(userKeystores, k)
		} else {
			for _, uid := range k.ViewerIds {
				if uid == userId {
					userKeystores = append(userKeystores, k)
				}
			}
		}
	}

	return userKeystores, nil
}

func (r *Repository) UserKeystore(userId, keystoreId string) (keystore.Keystore, error) {
	allKeystores, err := r.Keystores()
	if err != nil {
		return keystore.Keystore{}, err
	}

	for _, k := range allKeystores {
		if k.Id == keystoreId {
			if k.OwnerId == userId {
				return k, nil
			}

			for _, uid := range k.ViewerIds {
				if uid == userId {
					return k, nil
				}
			}
		}
	}

	return keystore.Keystore{}, keystore.ErrNotFound
}
