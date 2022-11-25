package mineservice

import (
	"io/ioutil"
	"log"

	"github.com/SKF/go-image-resizer/resizer"
)

func compressImage(imageName string) []byte {
	imageFile, err := ioutil.ReadFile(imageName)
	if err != nil {
		log.Fatal(err)
	}

	resizedImage, err := resizer.ResizeImage(imageFile, resizer.JpegEncoder, 1024, 1024)
	if err != nil {
		log.Fatal(err)
	}
	return resizedImage
}
