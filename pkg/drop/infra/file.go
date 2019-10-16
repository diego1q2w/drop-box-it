package infra

import (
	"fmt"
	"github.com/brainly/drop-box-it/pkg/drop/domain"
	"os"
	"path/filepath"
)

func ListFiles(root domain.Path) ([]domain.File, error) {
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
