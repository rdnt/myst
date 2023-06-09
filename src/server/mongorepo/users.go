package mongorepo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"myst/src/server/application"
	"myst/src/server/application/domain/user"
)

func (r *Repository) CreateUser(u user.User) (user.User, error) {
	collection := r.mdb.Database(r.database).Collection("users")
	ctx := context.Background()

	bu := UserToBSON(u)

	_, err := collection.InsertOne(ctx, bu)
	if err != nil {
		return user.User{}, errors.Wrap(err, "failed to insert user")
	}

	u = UserFromBSON(bu)

	return u, nil
}

func (r *Repository) User(id string) (user.User, error) {
	collection := r.mdb.Database(r.database).Collection("users")
	ctx := context.Background()

	res := collection.FindOne(ctx, bson.D{{"_id", id}})
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user.User{}, application.ErrUserNotFound
	} else if err != nil {
		return user.User{}, errors.Wrap(err, "failed to find user by id")
	}

	var u User
	err = res.Decode(&u)
	if err != nil {
		return user.User{}, errors.Wrap(err, "failed to decode user")
	}

	return UserFromBSON(u), nil
}

func (r *Repository) UserByUsername(username string) (user.User, error) {
	collection := r.mdb.Database(r.database).Collection("users")
	ctx := context.Background()

	res := collection.FindOne(ctx, bson.D{{"username", username}})
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user.User{}, application.ErrUserNotFound
	} else if err != nil {
		return user.User{}, errors.Wrap(err, "failed to find user by username")
	}

	var u User
	err = res.Decode(&u)
	if err != nil {
		return user.User{}, errors.Wrap(err, "failed to decode user")
	}

	return UserFromBSON(u), nil
}

func (r *Repository) Users() ([]user.User, error) {
	collection := r.mdb.Database(r.database).Collection("users")
	ctx := context.Background()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find users")
	}
	defer cur.Close(ctx)

	var u User
	users := []user.User{}
	for cur.Next(ctx) {
		err := cur.Decode(&u)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode user")
		}

		users = append(users, UserFromBSON(u))
	}
	if err := cur.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate over users")
	}

	return users, nil
}

func (r *Repository) UpdateUser(u user.User) (user.User, error) {
	collection := r.mdb.Database(r.database).Collection("users")
	ctx := context.Background()

	usr := UserToBSON(u)

	res := collection.FindOneAndReplace(ctx, bson.D{{"_id", u.Id}}, usr)
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user.User{}, application.ErrUserNotFound
	} else if err != nil {
		return user.User{}, errors.Wrap(err, "failed to find and replace user")
	}

	return UserFromBSON(usr), nil
}
