package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ctx         context.Context
	mongoclient *mongo.Client
)

func Connection() (db *mongo.Client) {
	return mongoclient
}

func ConnectToDB() *mongo.Client {
	logger := utility.NewLogger()
	uri := config.GetConfig().Mongodb.Url
	mongo_connection := options.Client().ApplyURI(uri)
	mongoClient, err := mongo.Connect(ctx, mongo_connection)
	if err != nil {
		log.Fatal(err)
	}

	// PINGING THE CONNECTION
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	// IF EVERYTHING IS OKAY, THEN CONNECT
	fmt.Println("MONGO CONNECTION ESTABLISHED")
	logger.Info("MONGO CONNECTION ESTABLISHED")

	mongoclient = mongoClient
	return mongoClient
}

// 1
// var Client *mongo.Client = ConnectToDB()

// getting database collections
func GetCollection(client *mongo.Client, databaseName, collectionName string) *mongo.Collection {
	collection := client.Database(databaseName).Collection(collectionName)
	return collection
}

func SelectFromCollection(ctx context.Context, database, collection string, filter bson.M) (*mongo.Cursor, error) {
	modelCollection := GetCollection(mongoclient, database, collection)
	cursor, err := modelCollection.Find(ctx, filter)
	if err != nil {
		return cursor, err
	}
	return cursor, nil
}
