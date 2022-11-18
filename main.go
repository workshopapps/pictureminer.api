package main

import (
<<<<<<< HEAD
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
=======
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/s3"
	"github.com/workshopapps/pictureminer.api/utility"

	// "github.com/workshopapps/pictureminer.api/pkg/repository/storage/redis"
	"github.com/workshopapps/pictureminer.api/pkg/router"
>>>>>>> e8168da4216a502374cedd40e23930a549d8ec23
)

func init() {
	config.Setup()
	// redis.SetupRedis() uncomment when you need redis
	mongodb.ConnectToDB()

<<<<<<< HEAD
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
=======
	s3.ConnectAws()
}

func main() {
	//Load config
	logger := utility.NewLogger()
	getConfig := config.GetConfig()
	validatorRef := validator.New()
	r := router.Setup(validatorRef, logger)

	logger.Info("Server is starting at 127.0.0.1:%s", getConfig.Server.Port)
	log.Fatal(r.Run(":" + getConfig.Server.Port))
>>>>>>> e8168da4216a502374cedd40e23930a549d8ec23
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
