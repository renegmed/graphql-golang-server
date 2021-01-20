package database

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client *mongo.Client
}

func NewDb() *DB {
	return connect()
}

var once sync.Once

func connect() *DB {

	var err error
	var client *mongo.Client

	once.Do(func() {
		client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			log.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = client.Connect(ctx)
	})
	if err != nil {
		log.Fatal(err)
	}

	return &DB{
		Client: client,
	}

}
