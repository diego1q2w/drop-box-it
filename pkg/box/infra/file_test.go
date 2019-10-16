package infra

import (
	"github.com/diego1q2w/drop-box-it/pkg/box/domain"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFile(t *testing.T) {
	rootDir := "/testdatabox"
	createDir(t, rootDir)
	file := domain.File{
		Path:    domain.Path(filepath.Join(rootDir, "test.txt")),
		Mode:    0755,
		Content: []byte(`heey this rocksa`),
	}

	fileBox := NewFileBox()
	err := fileBox.WriteFile(file)
	assert.NoError(t, err)

	content := readFileContent(t, file.Path.ToString())
	assert.Equal(t, content, file.Content)
}

func createDir(t *testing.T, path string) {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return
	}

	if err := os.MkdirAll(path, os.ModeDir); err != nil {
		t.Fatalf("unable to create test directory: %s", err)
	}
}

func readFileContent(t *testing.T, path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected error while reading: %s", err)
	}

	return data
}
