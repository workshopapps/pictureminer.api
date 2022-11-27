package mineservice

import (
	"io"
	"log"
	"mime/multipart"
	"strings"

	"github.com/h2non/bimg"
)

func compressImage(imageFile multipart.File) (io.ReadCloser, error) {
	buffer, err := io.ReadAll(imageFile)
	if err != nil {
		log.Fatal(err)
	}

	converted, err := bimg.NewImage(buffer).Convert(bimg.PNG)
	if err != nil {
		return nil, err
	}
	processed, err := bimg.NewImage(converted).Process(bimg.Options{Width: 1024, Height: 1024, Quality: 720})
	if err != nil {
		return nil, err
	}
	r := io.NopCloser(strings.NewReader(string(processed)))

	return r, nil
}
