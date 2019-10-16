package app

import (
	"context"
	"github.com/diego1q2w/drop-box-it/pkg/box/domain"
)

//go:generate moq -out file_boxer_mock_test.go . fileBoxer
type fileBoxer interface {
	WriteFile(file domain.File) error
	DeleteFile(path domain.Path) error
}

type Box struct {
	fileBoxer fileBoxer
}

func NewBox(fileBoxer fileBoxer) *Box {
	return &Box{fileBoxer: fileBoxer}
}

func (b *Box) WriteDocuments(ctx context.Context, file domain.File) error {
	return b.fileBoxer.WriteFile(file)
}

func (b *Box) DeleteDocuments(ctx context.Context, path domain.Path) error {
	return b.fileBoxer.DeleteFile(path)
}
