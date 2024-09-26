package cache

import (
	"os"
	"server/log"
)

var Files []*File

type File struct {
	Name    string
	Size    int
	Content *[]byte
	Handler *os.File
}

func LoadFiles(path []string) {

	for _, p := range path {
		NewFileCache(p)
	}

}

func generateFileCache(f *os.File) *File {
	info, _ := f.Stat()
	return &File{Name: info.Name(), Size: int(info.Size()), Handler: f}
}

func NewFileCache(p string) {

	dirs, err := os.ReadDir(p)

	if err != nil {
		if f, err := os.Open(p); err == nil {
			Files = append(Files, generateFileCache(f))
		} else {
			log.L.Sugar().Errorf("load file error: %s", err)
		}
		return
	}

	for _, dir := range dirs {

		if dir.IsDir() {
			NewFileCache(p + "/" + dir.Name())
			continue
		}
		f, err := os.Open(p + "/" + dir.Name())
		if err != nil {
			log.L.Sugar().Errorf("load file error: %s", err)
			continue
		}
		Files = append(Files, generateFileCache(f))

	}

}
