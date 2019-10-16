package infra

import (
	"fmt"
	"github.com/diego1q2w/drop-box-it/pkg/box/domain"
	"io/ioutil"
)

type FileBox struct {
}

func NewFileBox() *FileBox {
	return &FileBox{}
}

func (f *FileBox) WriteFile(file domain.File) error {
	if err := ioutil.WriteFile(file.Path.ToString(), file.Content, file.Mode); err != nil {
		return fmt.Errorf("unable to write document into a file")
	}

	return nil
}
