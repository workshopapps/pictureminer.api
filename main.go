package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ctx         context.Context
	mongoclient *mongo.Client
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	ctx = context.TODO()
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Please set MONGO_URI environment variable")
	}

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
}

func main() {
	server := gin.Default()
	defer mongoclient.Disconnect(ctx)

	// THIS IS THE BASEPATH TO THE ROUTES
	// basepath := server.Group("/v1")

	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "server online",
		})
	})

	// router.Run()
	port := os.Getenv("PORT")
	log.Fatal(server.Run(":" + port))
}
