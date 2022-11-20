package router

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
//
// 	"go.mongodb.org/mongo-driver/bson"
//
// 	"github.com/workshopapps/pictureminer.api/internal/model"
//
// 	"go.mongodb.org/mongo-driver/mongo"
//
// 	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
//
// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// )
//
// func setupRouter(client *mongo.Client) *gin.Engine {
// 	r := gin.Default()
// 	signupController := NewSignupController(client)
// 	r.POST("/signup", signupController.Handle)
//
// 	return r
// }
//
// func startMongoClient(t *testing.T) *mtest.T {
// 	options := mtest.NewOptions().
// 		ClientType(mtest.Mock).
// 		DatabaseName("ImageCollection").
// 		ShareClient(true)
//
// 	mt := mtest.New(t, options)
//
// 	mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "message", Value: "successful"}))
//
// 	return mt
// }
//
// func TestSignup_Handle_Unprocessable_Entity(t *testing.T) {
// 	mongoClient := startMongoClient(t)
// 	defer mongoClient.Close()
//
// 	router := setupRouter(mongoClient.Client)
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/signup", nil)
// 	router.ServeHTTP(w, req)
//
// 	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
// 	assert.Equal(t, `{"error":"could not read signup details"}`, w.Body.String())
// }
//
// func TestSignup_Handle_Successful(t *testing.T) {
// 	mongoClient := startMongoClient(t)
// 	defer mongoClient.Close()
//
// 	body, err := json.Marshal(model.UserStruct{
// 		UserName: "kdkkdkdkdk",
// 		Email:    "jackkd",
// 		Password: "©",
// 	})
//
// 	assert.NoError(t, err)
//
// 	router := setupRouter(mongoClient.Client)
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/signup", strings.NewReader(string(body)))
// 	router.ServeHTTP(w, req)
//
// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Equal(t, `{"userName":"kdkkdkdkdk","email":"jackkd"}`, w.Body.String())
// }
