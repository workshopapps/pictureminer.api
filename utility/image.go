package utility

import (
	"path/filepath"
	"strings"
)

func ValidImageFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg"
}
