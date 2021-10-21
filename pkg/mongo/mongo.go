package mongo

import (
	"context"

	"myst/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	log            = logger.New("mongo", logger.Blue)
	ErrNoDocuments = mongo.ErrNoDocuments
)

type Database struct {
	*mongo.Database
	client *mongo.Client
	name   string
	uri    string
}

func WithURI(uri string) func(*Database) {
	return func(db *Database) {
		db.uri = uri
	}
}

func WithName(name string) func(*Database) {
	return func(db *Database) {
		db.name = name
	}
}

func New(opts ...func(db *Database)) (*Database, error) {
	db := new(Database)
	for _, opt := range opts {
		opt(db)
	}

	ctx := context.Background()
	var err error
	db.client, err = mongo.Connect(ctx, options.Client().ApplyURI(db.uri))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	db.Database = db.client.Database(db.name)

	go func() {
		err = db.client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Error(err)
		} else {
			log.Print("MongoDB connected.")
		}
	}()

	return db, nil
}

func (db *Database) Close() {
	ctx := context.Background()
	err := db.client.Disconnect(ctx)
	if err != nil {
		log.Error(err)
	}
}

func DB() *mongo.Database {
	return nil
}
