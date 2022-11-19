package mongo

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
	var err error
	logger := utility.NewLogger()
	uri := config.GetConfig().Mongodb.Url
	mongo_connection := options.Client().ApplyURI(uri)
	mongoclient, err = mongo.Connect(ctx, mongo_connection)
	if err != nil {
		log.Fatal(err)
	}

	// PINGING THE CONNECTION
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	// IF EVERYTHING IS OKAY, THEN CONNECT
	fmt.Println("MONGO CONNECTION ESTABLISHED")
	logger.Info("MONGO CONNECTION ESTABLISHED")

	return mongoclient
}

func MongoPost(collection string, data interface{}) (interface{}, error) {
	c := getCollection(collection)

	result, err := c.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func MongoGet(collection string, filter map[string]interface{}) (interface{}, error) {
	c := getCollection(collection)

	f := make(bson.M, 0)

	if len(filter) == 1 {
		for k, v := range filter {
			f = bson.M{k: v}
		}
	} else if len(filter) > 1 {
		tf := make([]bson.M, 0, len(filter))
		for k, v := range filter {
			tf = append(tf, bson.M{k: v})
		}

		f = bson.M{"$and": tf}
	}

	cursor, err := c.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	results := make([]interface{}, 0)

	for cursor.Next(ctx) {
		var result interface{}
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func getCollection(collection string) *mongo.Collection {
	databaseName := config.GetConfig().Mongodb.Database
	database := mongoclient.Database(databaseName)
	c := database.Collection(collection)

	return c
}
