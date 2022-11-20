package microservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/model"
)

func GetImageContent(file io.ReadCloser, filename string) (*model.MicroserviceResponse, error) {
	microserviceHost := config.GetConfig().Python.MicroserviceHost

	req, err := SetupMultipartRequest(file, microserviceHost, "head.jpeg")
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Body)
	}

	var content model.MicroserviceResponse
	err = json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func SetupMultipartRequest(file io.ReadCloser, microserviceHost, filename string) (*http.Request, error) {
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", microserviceHost, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
