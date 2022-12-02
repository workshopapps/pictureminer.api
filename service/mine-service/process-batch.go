package mineservice

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/workshopapps/pictureminer.api/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	BatchName        = "name"
	BatchDescription = "description"
	BatchTags        = "tags"
	StatusOngoing    = "ongoing"
)

var (
	UrlMap = map[string]bool{
		"url":    true,
		"urls":   true,
		"image":  true,
		"images": true,
	}
)

func ProcessBatchService(file io.Reader) (interface{}, int, error) {
	// extract batch details and body from csv
	dMap, body, err := parseDetails(file)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// ensure all required details are available
	name, desc, tags, err := validateBatchDetails(dMap)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	urls, err := getURLs(body)
	fmt.Println(urls)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// run goroutine in background
	go processBatch(name, desc, tags, urls)

	// return success message
	res := model.BatchResponse{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Description: desc,
		Tags:        tags,
		Status:      StatusOngoing,
		DateCreated: time.Now().Local(),
	}
	return res, http.StatusOK, nil
}

func processBatch(name, desc string, tags, urls []string) {

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
	var res []string
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
