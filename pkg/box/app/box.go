package app

import "github.com/diego1q2w/drop-box-it/pkg/box/domain"

//go:generate moq -out file_boxer_mock_test.go . fileBoxer
type fileBoxer interface {
	WriteFile(file domain.File) error
}

type Box struct {
	fileBoxer fileBoxer
}

func NewBox(fileBoxer fileBoxer) *Box {
	return &Box{fileBoxer: fileBoxer}
}

func (b *Box) WriteDocuments(file domain.File) error {
	return b.fileBoxer.WriteFile(file)
}
