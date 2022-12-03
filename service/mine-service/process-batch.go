package mineservice

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	BatchName        = "name"
	BatchDescription = "description"
	BatchTags        = "tags"
	StatusOngoing    = "ongoing"
	StatusDone       = "done"
)

var (
	ERRNoURLsInJSON       = errors.New("no image urls specified in json file/ no csv file")
	ERRNoURLsInCSV        = errors.New("no urls in CSV")
	ERRURLFieldNotPresent = errors.New("url field not present in csv")
)

var (
	UrlMap = map[string]bool{
		"url":    true,
		"urls":   true,
		"image":  true,
		"images": true,
	}
)

func ProcessBatchCSVService(userID string, file io.Reader) (interface{}, int, error) {
	// extract batch details and body from csv
	dMap, body, err := parseDetails(file)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// ensure all required details are available
	bName, desc, tags, err := validateBatchDetails(dMap)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	urls, err := getURLs(body)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// save batch object to db
	batch := model.Batch{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		Name:        bName,
		Description: desc,
		Tags:        tags,
		Status:      StatusOngoing,
		DateCreated: time.Now().Local(),
	}
	database := config.GetConfig().Mongodb.Database
	batchCollection := mongodb.GetCollection(mongodb.Connection(), database, constants.BatchCollection)
	_, err = batchCollection.InsertOne(context.Background(), batch)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("%s: %s", "Unable to save user to database", err.Error())
	}

	uEmail, code, err := getUserEmail(userID)
	if err != nil {
		return nil, code, err
	}

	// run goroutine in background to process batch
	go processBatch(uEmail, bName, desc, batch.ID, tags, urls)

	// return success message
	res := model.BatchResponse{
		ID:          batch.ID,
		Name:        batch.Name,
		Description: batch.Description,
		Tags:        batch.Tags,
		Status:      batch.Status,
		DateCreated: batch.DateCreated,
	}
	return res, http.StatusOK, nil
}

func ProcessBatchService(userID, batchName, desc, tagsStr string, csvFile io.Reader) (interface{}, int, error) {
	// parse and filter valid tags
	tags := strings.Split(tagsStr, ",")
	validTags := []string{}
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			validTags = append(validTags, tag)
		}
	}

	// read urls from csv
	csvr := csv.NewReader(csvFile)
	dataset, err := csvr.ReadAll()
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// get index of urls
	idx := -1
	for i, header := range dataset[0] {
		if _, ok := UrlMap[strings.ToLower(header)]; ok {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil, http.StatusBadRequest, errors.New("could not find url column header")
	}

	// filter valid urls
	urls := []string{}
	for _, dataRow := range dataset[1:] {
		url := dataRow[idx]
		if isValidURL(url) {
			urls = append(urls, url)
		}
	}

	// save batch object to db
	batch := model.Batch{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		Name:        batchName,
		Description: desc,
		Tags:        validTags,
		Status:      StatusOngoing,
		DateCreated: time.Now().Local(),
	}
	database := config.GetConfig().Mongodb.Database
	batchCollection := mongodb.GetCollection(mongodb.Connection(), database, constants.BatchCollection)
	_, err = batchCollection.InsertOne(context.Background(), batch)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("%s: %s", "Unable to save user to database", err.Error())
	}

	uEmail, code, err := getUserEmail(userID)
	if err != nil {
		return nil, code, err
	}

	// run goroutine in background to process batch
	go processBatch(uEmail, batchName, desc, batch.ID, tags, urls)

	// return success message
	res := model.BatchResponse{
		ID:          batch.ID,
		Name:        batch.Name,
		Description: batch.Description,
		Tags:        batch.Tags,
		Status:      batch.Status,
		DateCreated: batch.DateCreated,
	}
	return res, http.StatusOK, nil
}

func parseDetails(file io.Reader) (map[string]string, []string, error) {
	details, body := getDetails(file)

	if details == nil || body == nil {
		return nil, nil, errors.New("invalid csv structure")
	}

	dMap := make(map[string]string)
	for _, line := range details {
		s := strings.Split(line, ":")
		if len(s) != 2 {
			return nil, nil, errors.New("invalid details structure for 'key:value' pair")
		}
		dMap[s[0]] = s[1]
	}

	return dMap, body, nil
}

func getDetails(file io.Reader) ([]string, []string) {
	details, body := []string{}, []string{}

	scn := bufio.NewScanner(file)
	for scn.Scan() {
		line := scn.Text()
		if strings.TrimSpace(line) == "" {
			break
		}
		details = append(details, line)
	}

	for scn.Scan() {
		line := scn.Text()
		if strings.TrimSpace(line) != "" {
			body = append(body, line)
		}
	}

	return details, body
}

func getURLs(body []string) ([]string, error) {
	headers := strings.Split(body[0], ",")

	// get index of header
	idx := -1
	for i, coln := range headers {
		if _, ok := UrlMap[strings.ToLower(coln)]; ok {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil, errors.New("could not find url column header")
	}

	// filter valid urls
	var res []string
	for _, row := range body[1:] {
		rs := strings.Split(row, ",")
		url := ""
		if idx < len(rs) {
			url = rs[idx]
		}
		if isValidURL(url) {
			res = append(res, url)
		}
	}

	return res, nil
}

func isValidURL(url string) bool {
	if url == "" {
		return false
	}

	ext := strings.ToLower(filepath.Ext(url))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg"
}

func validateBatchDetails(dMap map[string]string) (string, string, []string, error) {
	name, ok := dMap[BatchName]
	if !ok {
		return "", "", nil, errors.New("batch name missing from csv")
	}

	tag, ok := dMap[BatchTags]
	if !ok {
		return "", "", nil, errors.New("batch tags missing from csv")
	}

	// not required
	description := dMap[BatchDescription]

	tags := strings.Split(tag, ";")
	for i := 0; i < len(tags); i++ {
		tags[i] = strings.TrimSpace(tags[i])
	}

	return strings.TrimSpace(name), strings.TrimSpace(description), tags, nil
}

func getUserEmail(userID string) (string, int, error) {
	var user model.User

	database := config.GetConfig().Mongodb.Database
	userCollection := mongodb.GetCollection(mongodb.Connection(), database, constants.UserCollection)

	// convert "ObjectID('<id hex>') => '<id hex>'"
	id, err := primitive.ObjectIDFromHex(userID[10 : len(userID)-2])
	if err != nil {
		return "", http.StatusBadRequest, err
	}

	result := userCollection.FindOne(context.TODO(), bson.M{"_id": id})
	err = result.Err()
	if err != nil {
		return "", http.StatusNotFound, err
	}

	err = result.Decode(&user)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return user.Email, http.StatusOK, nil
}
