package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/workshopapps/pictureminer.api/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
)

var (
	ctx         context.Context
	mongoClient *mongo.Client
	db          *mongo.Database
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	ctx = context.TODO()
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Please set MO NGO_URI environment variable")
	}

	mongoConnection := options.Client().ApplyURI(uri)
	mongoClient, err = mongo.Connect(ctx, mongoConnection)
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
}

func main() {
	server := gin.Default()
	defer mongoClient.Disconnect(ctx)

	// THIS IS THE BASEPATH TO THE ROUTES
	// basepath := server.Group("/v1")

	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "server online",
		})
	})
	//SignUp route
	server.POST("/signup", SignUp)
	// router.Run()
	port := os.Getenv("PORT")
	log.Fatal(server.Run(":" + port))
}

func SignUp(c *gin.Context) {
	//Database connection
	imageDB := mongoClient.Database("ImageCollection")
	userCollection := imageDB.Collection("user")

	//Binding the userdetails to userStruct
	var User models.UserStruct
	err := c.Bind(&User)
	if err != nil {
		log.Fatal("Unable to bind user signup details")
	}
	//	Hashing the password
	harsh, err := bcrypt.GenerateFromPassword([]byte(User.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to harsh paswword",
		})
		return
	}
	User.Password = string(harsh)

	userCollection.InsertOne(context.Background(), User)
	c.JSON(200, User)
}
