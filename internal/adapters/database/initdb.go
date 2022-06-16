package database

import (
	"context"
	"fmt"
	"log"
	"medods/internal/adapters"
	"medods/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type db struct {
	collection *mongo.Collection
}

type Storage interface {
	GetOne(userGUID string, bindTokens string) (*model.UserToken, error)
	UpdateOne(old *model.UserToken, newToken *model.UserToken) error
	CreateOne(ut *model.UserToken) error
}

func NewCollection(mongo *mongo.Database, collection string) Storage {
	return &db{
		collection: mongo.Collection(collection),
	}
}

func NewDatabase(config adapters.Config) *mongo.Database {
	port := config.MongoURL
	fmt.Println(port)
	client, err := mongo.NewClient(options.Client().ApplyURI(port))
	if err != nil {
		log.Printf("Connection mongoDB: %v\n", err)
		log.Fatal()
	}
	if err = client.Connect(context.TODO()); err != nil {
		log.Printf("Client mongoDB: %v\n", err)
		log.Fatal()
	}
	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Printf("Client connection to mongoDB: %v\n", err)
		log.Fatal()
	}
	log.Println("Connection to MongoDB success")
	return client.Database(config.NameDB)
}
