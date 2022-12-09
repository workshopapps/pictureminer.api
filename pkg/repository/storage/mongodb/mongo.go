package mongodb

import (
	"time"
	"context"
	"fmt"
	"log"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ctx         context.Context
	mongoClient *mongo.Client
)

func Connection() (db *mongo.Client) {
	return mongoClient
}

func ConnectToDB() *mongo.Client {
	var err error
	logger := utility.NewLogger()
	uri := config.GetConfig().Mongodb.Url
	mongo_connection := options.Client().ApplyURI(uri)
	mongoClient, err = mongo.Connect(ctx, mongo_connection)
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

	return mongoClient
}

// 1
// var Client *mongo.Client = ConnectToDB()

// getting database collections
func GetCollection(client *mongo.Client, databaseName, collectionName string) *mongo.Collection {
	collection := client.Database(databaseName).Collection(collectionName)
	return collection
}

func MongoPost(collection string, data interface{}) (*mongo.InsertOneResult, error) {
	c := getCollection(collection)

	result, err := c.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func MongoGet(collection string, filter map[string]interface{}) ([]interface{}, error) {
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
	database := mongoClient.Database(databaseName)
	c := database.Collection(collection)

	return c
}

// Delete

func SelectFromCollection(ctx context.Context, database, collection string, filter bson.M) (*mongo.Cursor, error) {
	modelCollection := GetCollection(mongoClient, database, collection)
	cursor, err := modelCollection.Find(ctx, filter)
	if err != nil {
		return cursor, err
	}
	return cursor, nil
}

func DeleteAUserFromCollection(ctx context.Context, database, collection string, filter bson.M) (*mongo.DeleteResult, error) {
	modelCollection := GetCollection(mongoClient, database, collection)
	// To check if the user exist or not
	user := modelCollection.FindOne(ctx, filter)
	if user.Err() != nil {
		return nil, user.Err()
	}
	deletedResult, err := modelCollection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}
	return deletedResult, err

}

func CountFromCollection(user_id primitive.ObjectID) (int64, error) {
	userCollection := GetCollection(mongoClient, config.GetConfig().Mongodb.Database, constants.ImageCollection)
	filter := bson.D{{"user_id", user_id}}
	count, err := userCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return count, err
	}

	// implementaton code
	// estCount , CountErr := CountFromCollection("637f3cb921187c6a016f2087")
	// user.UsersCount = estCount
	// if CountErr != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error":"error reading number of documents"})
	// }
	return count, nil
}

func MainCountFromCollection(user_id string, collection string) (int64, error) {
	countCollection := GetCollection(mongoClient, config.GetConfig().Mongodb.Database, collection)
	filter := bson.D{{"user_id", user_id}}
	count, err := countCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return count, err
	}

	// implementaton code
	// count , err := MainCountFromCollection("637f3cb921187c6a016f2087")
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error":"error reading number of documents"})
	// }

	return count, nil
}


func GetUserTags(user_id string, batch_id primitive.ObjectID) ([]string, int, error) {
	var tags []string
	var length int

	batchImagesCollection := GetCollection(mongoClient, config.GetConfig().Mongodb.Database, constants.BatchCollection)
	filter := bson.D{{"user_id", user_id}, {"_id", batch_id}}

	batch_collection, err := batchImagesCollection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
	}

	var results []model.Batch
	if err := batch_collection.All(context.TODO(), &results); err != nil {
		fmt.Println(err)
	}

	for _, test := range results {
		if test.UserID == user_id {
			tags = test.Tags
			length = len(test.Tags)
		}
	}
	return tags, length, err
}

func GetImageTags(batch_id string) ([]model.BatchImage, []string, int, error) {
	var tag []string
	batchImagesCollection := GetCollection(mongoClient, config.GetConfig().Mongodb.Database, constants.BatchImageCollection)
	filter := bson.D{{"batch_id", batch_id}}

	image_collection, err := batchImagesCollection.Find(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}

	var results []model.BatchImage
	if err := image_collection.All(context.TODO(), &results); err != nil {
		fmt.Println(err)
	}

	for _, test := range results {
		if test.BatchID == batch_id {
			tag = append(tag, test.Tag)
		}
	}
	return results, tag, len(tag), err
}

func MongoUpdate(id string, updateEntries map[string]interface{}, collection string) (*mongo.UpdateResult, error) {
	c := getCollection(collection)
	apiCallCount := "api_call_count"
	user_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if value, ok := updateEntries[apiCallCount]; ok {
		_, err := c.UpdateByID(ctx, user_id, bson.M{"$inc": bson.M{apiCallCount: value}}, options.Update().SetUpsert(true))
		if err != nil {
			return nil, err
		}

		delete(updateEntries, apiCallCount)
		if len(updateEntries) == 0 {
			return nil, nil
		}
	}

	update := make(bson.M, 0)
	for i, j := range updateEntries {
		update[i] = j
	}

	db_data := bson.M{"$set": update}

	result, err := c.UpdateByID(context.TODO(), user_id, db_data)

	if err != nil {
		return nil, err
	}

	return result, nil
}



func GetUserPlan(user_id string) (string , error) {

	user_id_primitive, _ := primitive.ObjectIDFromHex(user_id)
	var plan string
	userCollection := GetCollection(mongoClient, config.GetConfig().Mongodb.Database, constants.UserCollection)
	filter := bson.D{{"user_id", user_id_primitive}}

	userCollectionPlan, err := userCollection.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println(err)
	}

	var results []model.User
	if err := userCollectionPlan.All(context.TODO(), &results); err != nil {
		fmt.Println(err)
	}

	for _, test := range results {
		if test.ID == user_id_primitive {
			plan = test.Plan
		}
	}
	return plan, err
}

func GetMinedTime(user_id string) ([]model.MinedImage, []time.Time, int, error) {
	var times []time.Time
	MinedImageinfo := GetCollection(mongoClient, config.GetConfig().Mongodb.Database, constants.ImageCollection)

	filter := bson.D{{"user_id", user_id}}

	image_time_collection, err := MinedImageinfo.Find(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}

	var results []model.MinedImage
	if err := image_time_collection.All(context.TODO(), &results); err != nil {
		fmt.Println(err)
	}

	for _, test := range results {
		localTime := time.Now().Local()
		toSubtract := -720 * time.Hour
		SubDate := localTime.Add(toSubtract)

		dateCreated := test.DateCreated

		if   SubDate.Before(dateCreated) {
			times = append(times, test.DateCreated)
		}
	}
	
	return results, times, len(times), err
}


func GetBatchTime(user_id string) ([]model.Batch, []time.Time, int, error) {
	var times []time.Time
	BatchImageinfo := GetCollection(mongoClient, config.GetConfig().Mongodb.Database, constants.BatchCollection)


	filter := bson.D{{"user_id", user_id}}
	// currentTime := time.Now()

	batch_time_collection, err := BatchImageinfo.Find(ctx, filter)
	if err != nil {
		fmt.Println(err)
	}

	var results []model.Batch
	if err := batch_time_collection.All(context.TODO(), &results); err != nil {
		fmt.Println(err)
	}

	for _, test := range results {
		localTime := time.Now().Local()
		toSubtract := -720 * time.Hour
		SubDate := localTime.Add(toSubtract)

		dateCreated := test.DateCreated

		if   SubDate.Before(dateCreated) {
			times = append(times, test.DateCreated)
		}
	}
	return results, times, len(times), err
}
