package mineservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/constants"
	"github.com/workshopapps/pictureminer.api/internal/model"
	"github.com/workshopapps/pictureminer.api/pkg/repository/storage/mongodb"
	"github.com/workshopapps/pictureminer.api/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Result struct {
	Name  string
	Score float64
}

type Label struct {
	URL     string
	Results []Result
}

type (
	APIResponse struct {
		Result APIResponseResult `json:"result"`
	}

	APIResponseResult struct {
		Tags []Tag `json:"tags"`
	}

	Tag struct {
		Confidence float64           `json:"confidence"`
		Tag        map[string]string `json:"tag"`
	}
)

const (
	Threshold = "30"
	Limit     = "10"
	Untagged  = "untagged"
)

func processBatch(email, batchID, name, desc string, tags, urls []string) {
	// labels for each url
	labels := fetchLabelsForURLS(urls)

	// classify label to matching tag
	batchImgs := classifyLabels(batchID, labels, tags)

	// save to db
	err := saveToDB(batchImgs)
	if err != nil {
		warnUser(email, name, err.Error())
	}

	// on complete, notify user through email
	notifyUser(email, name)
}

func fetchLabelsForURLS(urls []string) []Label {
	var labels []Label
	ImaggaURL := config.GetConfig().ImaggaAPI.URL
	AuthToken := config.GetConfig().ImaggaAPI.Auth
	httpClient := &http.Client{}
	for _, url := range urls {
		label, err := getLabel(httpClient, ImaggaURL, url, AuthToken, Limit, Threshold)
		if err == nil {
			labels = append(labels, label)g
		}
	}
	return labels
}

func getLabel(client *http.Client, imaggaURL, url, authToken, limit, threshold string) (Label, error) {
	req, _ := http.NewRequest("GET", imaggaURL, nil)
	req.Header.Set("Authorization", authToken)
	q := req.URL.Query()
	q.Add("limit", limit)
	q.Add("threshold", threshold)
	q.Add("image_url", url)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return Label{}, err
	}
	defer res.Body.Close()

	var apiRes APIResponse
	err = json.NewDecoder(res.Body).Decode(&apiRes)
	if err != nil {
		return Label{}, err
	}

	var results []Result
	for _, tag := range apiRes.Result.Tags {
		results = append(results, Result{Name: tag.Tag["en"], Score: tag.Confidence})
	}

	var label = Label{
		URL:     url,
		Results: results,
	}

	return label, nil
}

func classifyLabels(batchID string, labels []Label, tags []string) []model.BatchImage {
	// for O(1) lookups
	tagsMap := make(map[string]bool)
	for _, tag := range tags {
		tagsMap[smoothify(tag)] = true
	}

	// get best matching tag for each label
	var batchImgs []model.BatchImage
	for _, label := range labels {
		batchImgs = append(batchImgs, classifyLabel(batchID, label, tagsMap))
	}

	return batchImgs
}

func classifyLabel(batchID string, label Label, tagsMap map[string]bool) model.BatchImage {
	var filtered []Result

	// filter for results that match with atleast one tag
	for _, res := range label.Results {
		name := smoothify(res.Name)
		if tagsMap[name] {
			filtered = append(filtered, res)
		}
	}

	// get result with higest confidence/score
	bestTag, maxScore := Untagged, -1.0
	for _, res := range filtered {
		if res.Score > maxScore {
			bestTag, maxScore = res.Name, res.Score
		}
	}

	batchImage := model.BatchImage{
		ID:      primitive.NewObjectID(),
		BatchID: batchID,
		URL:     label.URL,
		Tag:     bestTag,
	}
	return batchImage
}

func smoothify(str string) string {
	return strings.Join(strings.Fields(str), "-")
}

func saveToDB(batchImgs []model.BatchImage) error {
	database := config.GetConfig().Mongodb.Database
	batchImgCol := mongodb.GetCollection(mongodb.Connection(), database, constants.BatchImageCollection)

	// mongodb insertMany supports []interface{} only
	imgs := make([]interface{}, len(batchImgs))
	for i := 0; i < len(imgs); i++ {
		imgs[i] = batchImgs[i]
	}

	_, err := batchImgCol.InsertMany(context.TODO(), imgs)
	if err != nil {
		return err
	}

	return nil
}

func notifyUser(email, batchName string) {
	from := config.GetConfig().NotifyEmail.Email
	password := config.GetConfig().NotifyEmail.Email
	body := fmt.Sprintf("Hello, batch <b>%v</b> processing is complete!", batchName)

	utility.EmailSender(from, password, []string{email}, "Process Batch Complete", body)
}

func warnUser(email, batchName, msg string) {
	from := config.GetConfig().NotifyEmail.Email
	password := config.GetConfig().NotifyEmail.Email
	body := fmt.Sprintf("Hello, batch <b>%v</b> processing is failed!<br>error: %v", batchName, msg)

	utility.EmailSender(from, password, []string{email}, "Process Batch Failed", body)
}
