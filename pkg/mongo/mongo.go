package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"myst/pkg/logger"
)

var (
	log = logger.New("mongo", logger.BlueFg)
	db  *mongo.Client
)

func New(uri string) (*mongo.Client, error) {
	ctx := context.Background()
	var err error
	db, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	go func() {
		err = db.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Error(err)
		} else {
			log.Print("MongoDB connected.")
		}
	}()
	return db, nil
}

// NewWithClient creates a new mongodb connection with the given mongodb client
func NewWithClient(mdb *mongo.Client) {
	db = mdb
}

func Close() {
	if db == nil {
		return
	}
	ctx := context.Background()
	err := db.Disconnect(ctx)
	if err != nil {
		log.Error(err)
	}
}
