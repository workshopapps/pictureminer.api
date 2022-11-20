package utility

import (
	"crypto/sha256"
	"fmt"
	"io"
)

func HashImage(image io.ReadCloser) (string, error) {

	hash := sha256.New()
	_, err := io.Copy(hash, image)
	if err != nil {
		return "", err
	}

	h := hash.Sum(nil)
	sha := fmt.Sprintf("%x", h)

	return sha, nil
}
