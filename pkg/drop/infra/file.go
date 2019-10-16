package infra

import (
	"fmt"
	"github.com/brainly/drop-box-it/pkg/drop/domain"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileFetcher struct {
}

func NewFileFetcher() *FileFetcher {
	return &FileFetcher{}
}

func (f *FileFetcher) ListFiles(root domain.Path) ([]domain.File, error) {
	var files []domain.File
	err := filepath.Walk(root.ToString(), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, domain.File{
				Path: domain.Path(path),
				Mode: info.Mode(),
			})
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to list files: %w", err)
	}

	return files, err
}

func (f *FileFetcher) ReadFileContent(path domain.Path) ([]byte, error) {
	data, err := ioutil.ReadFile(path.ToString())
	if err != nil {
		return nil, fmt.Errorf("unable to read file content: %w", err)
	}

	return data, nil
}
