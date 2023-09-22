package database

import (
	util "Patronus/util"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func ConnectDB() *mongo.Client {
	Mongo_URL := util.GoDotEnvVariable("MONGODB_URL")
	client, err := mongo.NewClient(options.Client().ApplyURI(Mongo_URL))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	defer cancel()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to mongoDB")
	return client
}

func GetUsersCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("PatronusDB").Collection("Users")
	return collection
}
