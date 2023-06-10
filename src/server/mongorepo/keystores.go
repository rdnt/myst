package mongorepo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"myst/src/server/application"
	"myst/src/server/application/domain/keystore"
)

func (r *Repository) CreateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	collection := r.mdb.Database(r.database).Collection("keystores")
	ctx := context.Background()

	kb := KeystoreToBSON(k)

	_, err := collection.InsertOne(ctx, kb)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to insert keystore")
	}

	k = KeystoreFromBSON(kb)

	return k, nil
}

func (r *Repository) Keystore(id string) (keystore.Keystore, error) {
	collection := r.mdb.Database(r.database).Collection("keystores")
	ctx := context.Background()

	res := collection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}})
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return keystore.Keystore{}, application.ErrKeystoreNotFound
	} else if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to find keystore")
	}

	var k Keystore
	err = res.Decode(&k)
	if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to decode keystore")
	}

	return KeystoreFromBSON(k), nil
}

func (r *Repository) Keystores() ([]keystore.Keystore, error) {
	collection := r.mdb.Database(r.database).Collection("keystores")
	ctx := context.Background()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find keystores")
	}
	defer cur.Close(ctx)

	var k Keystore
	keystores := []keystore.Keystore{}
	for cur.Next(ctx) {
		err := cur.Decode(&k)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode keystore")
		}

		keystores = append(keystores, KeystoreFromBSON(k))
	}
	if err := cur.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate over keystores")
	}

	return keystores, nil
}

func (r *Repository) UpdateKeystore(k keystore.Keystore) (keystore.Keystore, error) {
	collection := r.mdb.Database(r.database).Collection("keystores")
	ctx := context.Background()

	kb := KeystoreToBSON(k)
	res := collection.FindOneAndReplace(ctx, bson.D{primitive.E{Key: "_id", Value: kb.Id}}, kb)
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return keystore.Keystore{}, application.ErrKeystoreNotFound
	} else if err != nil {
		return keystore.Keystore{}, errors.Wrap(err, "failed to find and replace keystore")
	}

	return KeystoreFromBSON(kb), nil
}

func (r *Repository) DeleteKeystore(id string) error {
	collection := r.mdb.Database(r.database).Collection("keystores")
	ctx := context.Background()

	res := collection.FindOneAndDelete(ctx, bson.D{primitive.E{Key: "_id", Value: id}})
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return application.ErrKeystoreNotFound
	} else if err != nil {
		return errors.Wrap(err, "failed to find and delete keystore")
	}

	return nil
}
