package infra

import (
	"github.com/diego1q2w/drop-box-it/pkg/drop/domain"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFileFetcher_ListFiles(t *testing.T) {
	rootDir := "/tmp-drop-data"
	files := []domain.File{
		{Path: domain.Path(filepath.Join(rootDir, "lol", "test2.txt")), Mode: os.FileMode(0644)},
		{Path: domain.Path(filepath.Join(rootDir, "test1.txt")), Mode: os.FileMode(0755)},
		{Path: domain.Path(filepath.Join(rootDir, "test3.txt")), Mode: os.FileMode(0755)},
		{Path: domain.Path(filepath.Join(rootDir, "test4.txt")), Mode: os.FileMode(0755)},
	}
	removeDir(t, rootDir)
	createDir(t, rootDir)
	for _, file := range files {
		createFile(t, file, nil)
	}

	fileFetcher := NewFileFetcher()
	filesStored, err := fileFetcher.ListFiles(domain.Path(rootDir))
	assert.NoError(t, err)

	assert.Equal(t, files, filesStored)
}

func TestFileFetcher_IgnoreDirectories(t *testing.T) {
	rootDir := "/tmp-drop-data"
	nestedDir := filepath.Join(rootDir, "nested")
	files := []domain.File{
		{Path: domain.Path(filepath.Join(nestedDir, "test2.txt")), Mode: os.FileMode(0644)},
		{Path: domain.Path(filepath.Join(rootDir, "test1.txt")), Mode: os.FileMode(0755)},
	}
	removeDir(t, rootDir)
	createDir(t, nestedDir)
	for _, file := range files {
		createFile(t, file, nil)
	}

	fileFetcher := NewFileFetcher()
	filesStored, err := fileFetcher.ListFiles(domain.Path(rootDir))
	assert.NoError(t, err)

	assert.Equal(t, files, filesStored)
}

func TestReadFileContent(t *testing.T) {
	rootDir := "/tmp-drop-data"
	removeDir(t, rootDir)
	createDir(t, rootDir)
	filePath := domain.Path(filepath.Join(rootDir, "test.txt"))
	file := domain.File{Path: filePath, Mode: os.FileMode(0755)}
	content := []byte(`this should work`)

	createFile(t, file, content)

	fileFetcher := NewFileFetcher()
	storedContent, err := fileFetcher.ReadFileContent(filePath)
	assert.NoError(t, err)

	assert.Equal(t, content, storedContent)
}

func TestPathExists(t *testing.T) {
	rootDir := "/tmp-drop-data"
	removeDir(t, rootDir)
	createDir(t, rootDir)
	filePath := domain.Path(filepath.Join(rootDir, "test.txt"))
	file := domain.File{Path: filePath, Mode: os.FileMode(0755)}

	createFile(t, file, nil)
	fileFetcher := NewFileFetcher()

	exists, isDir := fileFetcher.PathExists(domain.Path(rootDir))
	assert.True(t, exists)
	assert.True(t, isDir)

	exists, isDir = fileFetcher.PathExists(filePath)
	assert.True(t, exists)
	assert.False(t, isDir)

	removeDir(t, rootDir)
	exists, isDir = fileFetcher.PathExists(filePath)
	assert.False(t, exists)
	assert.False(t, isDir)
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

func createFile(t *testing.T, file domain.File, data []byte) {
	dir := filepath.Dir(file.Path.ToString())
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		t.Fatalf("unable to create dir: %s", err)
	}

	if err := ioutil.WriteFile(file.Path.ToString(), data, file.Mode); err != nil {
		t.Fatalf("unable to create test file: %s", err)
	}
}

func removeDir(t *testing.T, path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return
	}

	if err := os.RemoveAll(path); err != nil {
		t.Fatalf("unable to delete test directory: %s", err)
	}
}
