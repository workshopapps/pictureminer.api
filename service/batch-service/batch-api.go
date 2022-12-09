package batchservice

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProcessBatchAPI(userID string, req *http.Request) (*model.ProcessBatchAPIResponse, error) {
	var urls []string
	var err error

	jsonReq, csvFile, err := processJSONAndCSV(req)
	if err != nil {
		return nil, err
	}

	if len(jsonReq.Images) != 0 {
		for _, v := range jsonReq.Images {
			if utility.ValidImageFormat(v) {
				urls = append(urls, v)
			}
		}
	} else {
		if csvFile == nil {
			return nil, ERRNoURLsInJSON
		}

		urls, err = parseCSVfile(csvFile)
		if err != nil {
			return nil, err
		}

		if len(urls) == 0 {
			return nil, ERRNoURLsInCSV
		}
	}

	time := time.Now()
	id, err := saveBatch(userID, jsonReq, time)
	if err != nil {
		return nil, err
	}

	userEmail, _, err := getUserEmail(userID)
	if err != nil {
		return nil, err
	}

	batchID := id.(primitive.ObjectID)

	jsonReq.Images = urls
	response := prepareResponse(batchID.Hex(), jsonReq, time)

	// start processing image urls in background
	go processBatch(userEmail, response.Name, response.Description, userID, batchID, response.Tags, jsonReq.Images)

	return response, nil
}

func saveBatch(uID string, pb *model.ProcessBatchAPIRequest, time time.Time) (interface{}, error) {

	batch := model.Batch{
		ID:          primitive.NewObjectID(),
		UserID:      uID,
		Name:        pb.Name,
		Description: pb.Description,
		Tags:        pb.Tags,
		Status:      StatusOngoing,
	}

	result, err := mongodb.MongoPost(constants.BatchCollection, batch)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

// processJSONFile processes the JSON file to get the json data and determines if the JSON file
// contains either image URLs or a CSV file is passed in with the request
func processJSONAndCSV(r *http.Request) (*model.ProcessBatchAPIRequest, multipart.File, error) {
	var jsonReq model.ProcessBatchAPIRequest

	file, _, err := r.FormFile("json")
	if err != nil {
		return nil, nil, err
	}

	if err = json.NewDecoder(file).Decode(&jsonReq); err != nil {
		return nil, nil, err
	}

	if len(jsonReq.Images) != 0 {
		return &jsonReq, nil, nil
	}

	csvFile, _, err := r.FormFile("csv")
	if err != nil {
		return nil, nil, err
	}

	return &jsonReq, csvFile, nil
}

func prepareResponse(id string, r *model.ProcessBatchAPIRequest, time time.Time) (response *model.ProcessBatchAPIResponse) {
	response = &model.ProcessBatchAPIResponse{
		ID:          id,
		Name:        r.Name,
		Description: r.Description,
		Tags:        r.Tags,
		Status:      StatusOngoing,
		DateCreated: time,
	}

	return
}
