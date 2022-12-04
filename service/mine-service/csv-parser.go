package mineservice

import (
	"encoding/csv"
	"mime/multipart"
	"strings"
	"os"

	"github.com/workshopapps/pictureminer.api/internal/model"
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

func ParseImageResponseForDownload(dt []model.BatchImage, ) error{
	//create the file
	file, err := os.Create("filename.csv")
	if err != nil{
		return err
	}
	writer := csv.NewWriter(file)

	var line []string
	l:=append(line, "url")
	l=append(l,"tag")
	writer.Write(l)
	//loop over the contents of the response
	for _, value := range dt{
		
		x :=append(line, value.URL)
		x =append(x, value.Tag)
		writer.Write(x)
	}
	writer.Flush()
	return nil
}