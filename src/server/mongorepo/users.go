package mongorepo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"myst/src/server/application/domain/user"
)

type User struct {
	Id           string    `bson:"_id"`
	Username     string    `bson:"username"`
	PasswordHash string    `bson:"passwordHash"`
	PublicKey    []byte    `bson:"publicKey"`
	CreatedAt    time.Time `bson:"createdAt"`
	UpdatedAt    time.Time `bson:"updatedAt"`
}

func UserToBSON(u user.User) User {
	return User{
		Id:           u.Id,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		PublicKey:    u.PublicKey,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func UserFromBSON(u User) user.User {
	return user.User{
		Id:           u.Id,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		PublicKey:    u.PublicKey,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func (r *Repository) CreateUser(u user.User) (user.User, error) {
	collection := r.db.Database("myst").Collection("users")

	bu := UserToBSON(u)

	_, err := collection.InsertOne(context.Background(), bu)
	if err != nil {
		return user.User{}, err
	}

	u = UserFromBSON(bu)

	return u, nil
}

func (r *Repository) User(id string) (user.User, error) {
	collection := r.db.Database("myst").Collection("users")

	res := collection.FindOne(context.Background(), bson.D{{"_id", id}})
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user.User{}, user.ErrNotFound
	} else if err != nil {
		return user.User{}, err
	}

	var u User
	err = res.Decode(&u)
	if err != nil {
		return user.User{}, err
	}

	return UserFromBSON(u), nil
}

func (r *Repository) UserByUsername(username string) (user.User, error) {
	collection := r.db.Database("myst").Collection("users")

	res := collection.FindOne(context.Background(), bson.D{{"username", username}})
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user.User{}, user.ErrNotFound
	} else if err != nil {
		return user.User{}, err
	}

	var u User
	err = res.Decode(&u)
	if err != nil {
		return user.User{}, err
	}

	return UserFromBSON(u), nil
}

func (r *Repository) Users() ([]user.User, error) {
	collection := r.db.Database("myst").Collection("users")

	ctx := context.Background()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var u User
	users := []user.User{}
	for cur.Next(ctx) {
		err := cur.Decode(&u)
		if err != nil {
			return nil, err
		}

		users = append(users, UserFromBSON(u))
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) UpdateUser(u user.User) (user.User, error) {
	collection := r.db.Database("myst").Collection("users")

	ctx := context.Background()

	usr := UserToBSON(u)
	res := collection.FindOneAndReplace(ctx, bson.D{{"_id", u.Id}}, usr)
	err := res.Err()
	if errors.Is(err, mongo.ErrNoDocuments) {
		return user.User{}, user.ErrNotFound
	} else if err != nil {
		return user.User{}, err
	}

	return UserFromBSON(usr), nil
}
