package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"myst/logger"
)

var (
	log         = logger.New("database", logger.BlueFg)
	client      *mongo.Client
	db          *mongo.Database
	ErrNotReady = fmt.Errorf("not ready")
)

func New(uri string) (*mongo.Database, error) {
	ctx := context.Background()
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	go func() {
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Error(err)
		} else {
			log.Print("MongoDB connected.")
		}
	}()
	db = client.Database("myst")
	return db, nil
}

// NewWithClient creates a new mongodb connection with the given mongodb client
func NewWithClient(c *mongo.Client) *mongo.Database {
	client = c
	db = client.Database("myst")
	return db
}

func DB() *mongo.Database {
	return db
}

func Close() {
	if db == nil {
		return
	}
	ctx := context.Background()
	err := client.Disconnect(ctx)
	if err != nil {
		log.Error(err)
	}
}
