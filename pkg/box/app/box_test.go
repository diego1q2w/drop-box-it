package app

import (
	"context"
	"fmt"
	"github.com/diego1q2w/drop-box-it/pkg/box/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileBoxerWriteError(t *testing.T) {
	expectedError := fmt.Errorf("test")
	expectedFile := domain.File{
		Path: domain.Path("root/test.txt"),
	}
	fileBoxer := &fileBoxerMock{
		WriteFileFunc: func(file domain.File) error {
			assert.Equal(t, expectedFile, file)
			return fmt.Errorf("test")
		},
	}

	box := NewBox(fileBoxer, "./root")
	err := box.WriteDocuments(context.Background(), domain.File{Path: "test.txt"})
	if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", expectedError) {
		t.Errorf("exepected error: %s, got: %s", expectedError, err)
	}
}

func TestFileBoxerDeleteError(t *testing.T) {
	expectedError := fmt.Errorf("test")
	expectedPath := domain.Path("root/test.txt")
	fileBoxer := &fileBoxerMock{
		DeleteFileFunc: func(path domain.Path) error {
			assert.Equal(t, expectedPath, path)
			return fmt.Errorf("test")
		},
	}

	box := NewBox(fileBoxer, "./root")
	err := box.DeleteDocuments(context.Background(), "test.txt")
	if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", expectedError) {
		t.Errorf("exepected error: %s, got: %s", expectedError, err)
	}
}
