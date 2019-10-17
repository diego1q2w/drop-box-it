package app

import (
	"context"
	"github.com/diego1q2w/drop-box-it/pkg/box/domain"
	"log"
	"path/filepath"
)

//go:generate moq -out file_boxer_mock_test.go . fileBoxer
type fileBoxer interface {
	WriteFile(file domain.File) error
	DeleteFile(path domain.Path) error
}

type Box struct {
	fileBoxer fileBoxer
	rootPath  string
}

func NewBox(fileBoxer fileBoxer, rootPath string) *Box {
	return &Box{fileBoxer: fileBoxer, rootPath: rootPath}
}

func (b *Box) WriteDocuments(ctx context.Context, file domain.File) error {
	log.Printf("writting document with path '%s'", file.Path)
	file.Path = domain.Path(filepath.Join(b.rootPath, file.Path.ToString()))
	return b.fileBoxer.WriteFile(file)
}

func (b *Box) DeleteDocuments(ctx context.Context, path domain.Path) error {
	log.Printf("deleting document with path '%s'", path)
	path = domain.Path(filepath.Join(b.rootPath, path.ToString()))
	return b.fileBoxer.DeleteFile(path)
}
