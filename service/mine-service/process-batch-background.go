package mineservice

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/workshopapps/pictureminer.api/internal/config"
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
	Threshold = "40"
	Limit     = "5"
)

func processBatch(email, batchID, name, desc string, tags, urls []string) {
	// labels for each url
	labels := fetchLabelsForURLS(urls)

	for _, label := range labels {
		fmt.Println(label)
	}

	// classify url to tag

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

func classifyLabels() {

}

func saveToDB() {

}

func nofityUser() {

}
