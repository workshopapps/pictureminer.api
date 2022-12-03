package batchservice

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"
)

// returns the index of the url
func getUrlHeaderIndex(headers []string) (int, error) {
	valid_headers := []string{"url", "urls", "image", "images"}

	for i, head := range headers {
		if head == valid_headers[0] || head == valid_headers[1] || head == valid_headers[2] || head == valid_headers[3] {
			return i, nil
		}
	}
	return 0, errors.New("no valid csv header present")
}

func checkExtension(url string) bool {
	mime_types := []string{".png", ".jpg", ".jpeg"}

	for _, mime_type := range mime_types {
		if strings.HasSuffix(url, mime_type) {
			return true
		}
	}
	return false
}

func parseCSVfile(file io.Reader) ([]string, error) {

	var urls []string // slice to store urls

	// Read CSV file
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return []string{}, err
	}

	var headers []string // csv headers
	for row, column := range records {
		if row == 0 {
			for _, v := range column {
				headers = append(headers, v)
			}
		}
	}

	index, err := getUrlHeaderIndex(headers)
	if err != nil {
		return []string{}, err
	}

	for row, column := range records {
		// Skip the first row of the csv file which contains the header
		if row == 0 {
			continue
		}

		url := column[index]

		// check if url is blank, then skip
		if url == "" {
			continue
		}

		// if it returns true, append all the valid url to urls
		if checkExtension(url) {
			urls = append(urls, url)
		}

	}

	return urls, nil
}
