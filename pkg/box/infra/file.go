package infra

import (
	"fmt"
	"github.com/diego1q2w/drop-box-it/pkg/box/domain"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileBox struct {
}

func NewFileBox() *FileBox {
	return &FileBox{}
}

func (f *FileBox) WriteFile(file domain.File) error {
	dir := filepath.Dir(file.Path.ToString())
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("unable to create dir: %w", err)
	}

	if err := ioutil.WriteFile(file.Path.ToString(), file.Content, file.Mode); err != nil {
		return fmt.Errorf("unable to write document into a file: %w", err)
	}

	return nil
}

func (f *FileBox) DeleteFile(path domain.Path) error {
	if err := os.Remove(path.ToString()); err != nil {
		return fmt.Errorf("unable to delete file")
	}

	return nil
}
