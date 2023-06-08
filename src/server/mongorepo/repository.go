package mongorepo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	mdb      *mongo.Client
	database string
}

func (r *Repository) DropDatabase() error {
	err := r.mdb.Database(r.database).Drop(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to drop database")
	}

	return nil
}

func New(addr string, database string) (*Repository, error) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(addr))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to mongodb")
	}

	return &Repository{
		mdb:      client,
		database: database,
	}, nil
}
