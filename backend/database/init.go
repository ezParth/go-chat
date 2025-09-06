package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect() {
	MONGO_URI := "mongodb://localhost:27017"
	clientOption := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		log.Fatal("Error in mongo.Connect -> ", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error in client.Ping: ", err)
	}

	// userCollection = client.Database("go-chat").Collection("User")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("Connected to db")

	Client = client

	// return userCollection
	// return client
}
