package cache

import (
	"io"
	"os"
	path2 "path"
)

var Files []*File

type File struct {
	Name    string
	Size    int
	Content *[]byte
}

func LoadFiles(path string) error {
	cache, err := NewFileCache(path)
	if err != nil {
		return err
	}
	Files = append(Files, cache)
	return nil
}

func NewFileCache(path string) (*File, error) {

	open, err := os.Open(path)
	_, file := path2.Split(path)
	defer open.Close()
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(open)

	if err != nil {
		return nil, err
	}

	return &File{Name: file, Size: len(all), Content: &all}, nil
}
