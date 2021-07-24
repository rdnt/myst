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
)

func New(uri string) (*mongo.Client, error) {
	ctx := context.Background()
	db, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
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

func Close(db *mongo.Client) {
	if db == nil {
		return
	}
	ctx := context.Background()
	err := db.Disconnect(ctx)
	if err != nil {
		log.Error(err)
	}
}
