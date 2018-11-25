package controllers

import (
	"io/ioutil"
	"path"
	"time"
)

type File struct {
	Name        string    `json:"name"`
	IsDirectory bool      `json:"isDirectory"`
	Items       []File    `json:"items"`
	Size        int64     `json:"size"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func getDirFiles(dirPath string) ([]File, error) {
	var err error

	files := make([]File, 0)

	ff, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, f := range ff {
		var items []File

		if f.IsDir() {
			items, err = getDirFiles(path.Join(dirPath, f.Name()))
			if err != nil {
				return nil, err
			}
		}

		files = append(files, File{
			Name:        f.Name(),
			IsDirectory: f.IsDir(),
			Items:       items,
			Size:        f.Size(),
			UpdatedAt:   f.ModTime(),
		})
	}

	return files, nil
}
