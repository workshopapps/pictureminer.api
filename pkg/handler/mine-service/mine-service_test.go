package mineservice

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/s3"
	"github.com/workshopapps/pictureminer.api/utility"
)

func setUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}
func init() {
	config.Setup()
	mongodb.ConnectToDB()
	s3.ConnectAws()
}

func TestMineServiceAPI(t *testing.T) {
	logger := utility.NewLogger()
	validatorRef := validator.New()

	r := setUpRouter()

	filePath := "image.png"
	fieldName := "image"
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	w, err := writer.CreateFormFile(fieldName, filePath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(w, file); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		log.Fatal(err)
	}
	Mine := Controller{validatorRef, logger}

	r.POST("/api/v1/mine-service", Mine.Post)

	req := httptest.NewRequest("POST", "/api/v1/mine-service", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res := httptest.NewRecorder()

	r.ServeHTTP(res, req)
	response := res.Result()
	resp, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(resp))

	assert.Equal(t, http.StatusOK, res.Code)

	t.Log("It should respond with an HTTP status code of 200")
	if res.Code != 200 {
		t.Errorf("Expected %d, received %d", 200, res.Code)
	}
}
