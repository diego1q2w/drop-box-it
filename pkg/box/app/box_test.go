package app

import (
	"context"
	"github.com/diego1q2w/drop-box-it/pkg/box/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFileBoxerWriteWasCalled(t *testing.T) {
	expectedFile := domain.File{
		Path: domain.Path("root/test.txt"),
	}
	fileBoxer := &fileBoxerMock{
		WriteFileFunc: func(file domain.File) error {
			return nil
		},
	}

	hasher := func(s string) uint64 {
		return 5
	}
	box := NewBox(fileBoxer, hasher, "./root", 5)
	err := box.WriteDocuments(context.Background(), domain.File{Path: "test.txt"})
	assert.NoError(t, err)
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, 1, len(fileBoxer.calls.WriteFile))
	assert.Equal(t, expectedFile, fileBoxer.calls.WriteFile[0].File)
}

func TestFileBoxerDeleteWasCalled(t *testing.T) {
	expectedPath := domain.Path("root/test.txt")
	fileBoxer := &fileBoxerMock{
		DeleteFileFunc: func(path domain.Path) error {
			return nil
		},
	}

	hasher := func(s string) uint64 {
		return 7
	}
	box := NewBox(fileBoxer, hasher, "./root", 5)
	err := box.DeleteDocuments(context.Background(), "test.txt")
	assert.NoError(t, err)
	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, 1, len(fileBoxer.calls.DeleteFile))
	assert.Equal(t, expectedPath, fileBoxer.calls.DeleteFile[0].Path)
}
