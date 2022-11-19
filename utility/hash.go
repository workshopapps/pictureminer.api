package utility

import (
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
)

func HashImage(image multipart.File) (string, error) {

	hash := sha256.New()
	_, err := io.Copy(hash, image)
	if err != nil {
		return "", err
	}

	h := hash.Sum(nil)
	sha := fmt.Sprintf("%x", h)

	return sha, nil
}
