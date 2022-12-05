package mineservice

import (
	"encoding/csv"
	"mime/multipart"
	"strings"
)

func checkExtension(url string) bool {
	mime_types := []string{".png", ".jpg", ".jpeg"}

	for _, mime_type := range mime_types {
		if strings.HasSuffix(url, mime_type) {
			return true
		}
	}
	return false
}

// func ParseCSVfile(file *os.File) []string {
func ParseCSVfile(file multipart.File) ([]string, error) {

	var urls []string

	// Read CSV file
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return []string{}, err
	}

	for row, column := range records {
		// Skip the first row of the csv file which contains the header
		if row == 0 {
			continue
		}

		url := column[0]
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
