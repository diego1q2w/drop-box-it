package app

import (
	"context"
	"fmt"
	"github.com/diego1q2w/drop-box-it/pkg/drop/domain"
	"sync"
)

type isDir = bool
type exists = bool

//go:generate moq -out file_fetcher_mock_test.go . fileFetcher
type fileFetcher interface {
	ListFiles(root domain.Path) ([]domain.File, error)
	ReadFileContent(path domain.Path) ([]byte, error)
	PathExists(path domain.Path) (exists, isDir)
}

type Dropper struct {
	fileFetcher fileFetcher
	mux         sync.Mutex
	filesStatus map[domain.Path]domain.FileStatus
}

func NewDropper(fetcher fileFetcher) *Dropper {
	return &Dropper{fileFetcher: fetcher}
}

func (d *Dropper) SyncFiles(ctx context.Context, rootPath domain.Path) error {
	if err := d.validateRootPath(rootPath); err != nil {
		return err
	}

	files, err := d.fileFetcher.ListFiles(rootPath)
	if err != nil {
		return fmt.Errorf("error while listing files: %w", err)
	}

	d.updateFileStatuses(files)
	return nil
}

func (d *Dropper) updateFileStatuses(files []domain.File) {
	d.mux.Lock()
	defer d.mux.Unlock()

	for _, file := range files {
		if _, ok := d.filesStatus[file.Path]; ok {
			d.filesStatus[file.Path] = domain.Updated
		} else {
			d.filesStatus[file.Path] = domain.Created
		}
	}

	for path := range d.filesStatus {
		if !containsPath(files, path) {
			d.filesStatus[path] = domain.Deleted
		}
	}
}

func containsPath(s []domain.File, e domain.Path) bool {
	for _, a := range s {
		if a.Path == e {
			return true
		}
	}
	return false
}

func (d *Dropper) validateRootPath(rootPath domain.Path) error {
	exists, isDir := d.fileFetcher.PathExists(rootPath)
	if !exists {
		return fmt.Errorf("the path '%s' does not exists", rootPath)
	}

	if !isDir {
		return fmt.Errorf("the path priovided is not a directory")
	}

	return nil
}
