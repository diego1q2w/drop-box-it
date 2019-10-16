package app

import (
	"fmt"
	"github.com/diego1q2w/drop-box-it/pkg/box/domain"
	"testing"
)

func TestFileBoxerWriteError(t *testing.T) {
	expectedError := fmt.Errorf("test")
	fileBoxer := &fileBoxerMock{
		WriteFileFunc: func(file domain.File) error {
			return fmt.Errorf("test")
		},
	}

	box := NewBox(fileBoxer)
	err := box.WriteDocuments(domain.File{})
	if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", expectedError) {
		t.Errorf("exepected error: %s, got: %s", expectedError, err)
	}
}

func TestFileBoxerDeleteError(t *testing.T) {
	expectedError := fmt.Errorf("test")
	fileBoxer := &fileBoxerMock{
		DeleteFileFunc: func(path domain.Path) error {
			return fmt.Errorf("test")
		},
	}

	box := NewBox(fileBoxer)
	err := box.DeleteDocuments(domain.File{})
	if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", expectedError) {
		t.Errorf("exepected error: %s, got: %s", expectedError, err)
	}
}
