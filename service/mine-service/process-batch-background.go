package mineservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/workshopapps/pictureminer.api/internal/config"
	"github.com/workshopapps/pictureminer.api/internal/model"
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

var (
	BatchId = ""
)

func processBatch(email, batchID, name, desc string, tags, urls []string) {
	BatchId = batchID

	// labels for each url
	labels := fetchLabelsForURLS(urls)

	// classify label to matching tag
	batchImgs := classifyLabels(labels, tags)
	for _, bImg := range batchImgs {
		fmt.Println(bImg)
	}

	// save to db

	// on complete, notify user through email
}

func fetchLabelsForURLS(urls []string) []Label {
	var labels []Label
	ImaggaURL := config.GetConfig().ImaggaAPI.URL
	AuthToken := config.GetConfig().ImaggaAPI.Auth
	httpClient := &http.Client{}
	for _, url := range urls {
		label := getLabel(httpClient, ImaggaURL, url, AuthToken, Limit, Threshold)
		labels = append(labels, label)
	}
	return labels
}

func getLabel(client *http.Client, imaggaURL, url, authToken, limit, threshold string) Label {
	req, _ := http.NewRequest("GET", imaggaURL, nil)
	req.Header.Set("Authorization", authToken)
	q := req.URL.Query()
	q.Add("limit", limit)
	q.Add("threshold", threshold)
	q.Add("image_url", url)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	var apiRes APIResponse
	err = json.NewDecoder(res.Body).Decode(&apiRes)
	if err != nil {
		fmt.Println(err)
	}

	var results []Result
	for _, tag := range apiRes.Result.Tags {
		results = append(results, Result{Name: tag.Tag["en"], Score: tag.Confidence})
	}

	var label = Label{
		URL:     url,
		Results: results,
	}

	return label
}

func classifyLabels(labels []Label, tags []string) []model.BatchImage {
	// for O(1) lookups
	tagsMap := make(map[string]bool)
	for _, tag := range tags {
		tagsMap[smoothify(tag)] = true
	}

	// get best matching tag for each label
	var batchImgs []model.BatchImage
	for _, label := range labels {
		batchImgs = append(batchImgs, classifyLabel(label, tagsMap))
	}

	return batchImgs
}

func classifyLabel(label Label, tagsMap map[string]bool) model.BatchImage {
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
		BatchID: BatchId,
		URL:     label.URL,
		Tag:     bestTag,
	}
	return batchImage
}

func smoothify(str string) string {
	return strings.Join(strings.Fields(str), "-")
}

func saveToDB() {

}

func nofityUser() {

}
