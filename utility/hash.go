package utility

import (
	"crypto/sha256"
	"fmt"
	"io"
)

func HashFile(file io.ReadCloser) (string, error) {

	hash := sha256.New()
	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
