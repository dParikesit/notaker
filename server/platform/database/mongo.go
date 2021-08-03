package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const mongoURI = "mongodb://localhost:27017/"

var Client *mongo.Client
var db *mongo.Database
var Notes *mongo.Collection

func Connect() error {
	var err error
	Client, err = mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = Client.Connect(ctx)
	if err != nil {
		log.Panic(err)
	}

	err = Client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Panic(err)
	}

	db = Client.Database("notaker")
	Notes = db.Collection("notes")

	return nil
}
