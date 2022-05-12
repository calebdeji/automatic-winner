package database

import (
	"context"
	"log"
	"os"
	"zate/services"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {

	services.LoadEnv()

	dbURL := os.Getenv("MONGODB_URL")

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURL))

	if err != nil {
		log.Printf(err.Error())
		log.Fatal("Unable to connect to mongo")

	}

	// defer disconnectClient(client, ctx)

	return client

}

func disconnectClient(client *mongo.Client, ctx context.Context) {
	if err := client.Disconnect(ctx); err != nil {
		log.Panic(err)
	}
}

var Client *mongo.Client = DBInstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	dbName := os.Getenv("DBNAME")

	var collection *mongo.Collection = client.Database(dbName).Collection(collectionName)
	return collection
}
