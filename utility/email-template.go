package utility

import (
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
)

// creating function to dynamically parse email template
func EmailTemplateParser(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	fmt.Println("Parsing the Template..")

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}