package mongorepo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	mdb      *mongo.Client
	database string
}

func (r *Repository) FlushDB() error {
	return r.mdb.Database(r.database).Drop(context.Background())
}

func New(addr string, database string) (*Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(addr))
	if err != nil {
		return nil, err
	}

	return &Repository{
		mdb:      client,
		database: database,
	}, nil
}
