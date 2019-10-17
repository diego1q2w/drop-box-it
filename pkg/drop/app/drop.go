package app

import (
	"context"
	"fmt"
	"github.com/diego1q2w/drop-box-it/pkg/drop/domain"
	"log"
	"path/filepath"
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

//go:generate moq -out box_client_mock_test.go . boxClient
type boxClient interface {
	WriteDocument(ctx context.Context, file domain.File, content []byte) error
	DeleteDocument(ctx context.Context, file domain.File, content []byte) error
}

type fileStatus struct {
	status domain.FileStatus
	file   domain.File
}

type Dropper struct {
	fileFetcher fileFetcher
	boxClient   boxClient
	mux         sync.Mutex
	rootPath    domain.Path
	filesStatus map[domain.Path]fileStatus
}

func NewDropper(fetcher fileFetcher, boxClient boxClient, rootPath string) *Dropper {
	return &Dropper{
		fileFetcher: fetcher,
		boxClient:   boxClient,
		rootPath:    domain.Path(rootPath),
		filesStatus: make(map[domain.Path]fileStatus),
	}
}

func (d *Dropper) SyncFiles(ctx context.Context) error {
	if err := d.validateRootPath(); err != nil {
		return err
	}

	files, err := d.fileFetcher.ListFiles(d.rootPath)
	if err != nil {
		return fmt.Errorf("error while listing files: %w", err)
	}

	files = d.fetchContent(files)

	d.updateFileStatuses(files)
	return nil
}

func (d *Dropper) fetchContent(files []domain.File) []domain.File {
	var fileWithContent []domain.File
	for _, file := range files {
		content, err := d.fileFetcher.ReadFileContent(file.Path.WithRoot(d.rootPath))
		if err != nil {
			log.Printf("unable to fetch file content: %s", err)
			continue
		}

		file.Content = domain.Content(content).Hash()
		fileWithContent = append(fileWithContent, file)
	}

	return fileWithContent
}

func (d *Dropper) updateFileStatuses(files []domain.File) {
	d.mux.Lock()
	defer d.mux.Unlock()

	for _, file := range files {
		file.Path = file.Path.RemoveBasePath(d.rootPath)
		if _, ok := d.filesStatus[file.Path]; ok {
			if file.Equal(d.filesStatus[file.Path].file) {
				continue
			}

			d.filesStatus[file.Path] = fileStatus{
				status: domain.Updated,
				file:   file,
			}
		} else {
			d.filesStatus[file.Path] = fileStatus{
				status: domain.Created,
				file:   file,
			}
		}
	}

	for path := range d.filesStatus {
		if !d.containsPath(files, path) {
			fileStatus := d.filesStatus[path]
			fileStatus.status = domain.Deleted
			d.filesStatus[path] = fileStatus
		}
	}
}

func (d *Dropper) containsPath(files []domain.File, path domain.Path) bool {
	for _, file := range files {
		file.Path = file.Path.RemoveBasePath(d.rootPath)
		if file.Path == path {
			return true
		}
	}
	return false
}

func (d *Dropper) SendUpdates(ctx context.Context) {
	paths := d.getPathsByStatus(domain.Created)
	d.writeDocuments(ctx, paths)

	paths = d.getPathsByStatus(domain.Updated)
	d.writeDocuments(ctx, paths)

	paths = d.getPathsByStatus(domain.Deleted)
	d.deleteDocuments(ctx, paths)
}

func (d *Dropper) writeDocuments(ctx context.Context, files []domain.File) {
	var wp sync.WaitGroup
	for _, file := range files {
		content, err := d.fileFetcher.ReadFileContent(file.Path.WithRoot(d.rootPath))
		if err != nil {
			log.Printf("unable to fetch file content: %s", err)
			return
		}

		wp.Add(1)
		go func(file domain.File) {
			defer wp.Done()
			d.writeDocument(ctx, file, content)
		}(file)
	}
	wp.Wait()
}

func (d *Dropper) writeDocument(ctx context.Context, file domain.File, content []byte) {
	if err := d.boxClient.WriteDocument(ctx, file, content); err != nil {
		log.Printf("error while writing document: %s", err)
		return
	}

	d.mux.Lock()
	defer d.mux.Unlock()
	fileStatus := d.filesStatus[file.Path]
	fileStatus.status = domain.Synced
	d.filesStatus[file.Path] = fileStatus
}

func (d *Dropper) deleteDocuments(ctx context.Context, files []domain.File) {
	var wp sync.WaitGroup
	for _, file := range files {
		wp.Add(1)
		go func(file domain.File) {
			defer wp.Done()
			d.deleteDocument(ctx, file)
		}(file)
	}
	wp.Wait()
}

func (d *Dropper) deleteDocument(ctx context.Context, file domain.File) {
	if err := d.boxClient.DeleteDocument(ctx, file, nil); err != nil {
		log.Printf("error while deleting document: %s", err)
		return
	}

	d.mux.Lock()
	defer d.mux.Unlock()
	fileStatus := d.filesStatus[file.Path]
	if fileStatus.status == domain.Deleted {
		delete(d.filesStatus, file.Path)
	}
}

func (d *Dropper) getPathsByStatus(fileStatus domain.FileStatus) []domain.File {
	d.mux.Lock()
	defer d.mux.Unlock()

	var files []domain.File
	for _, status := range d.filesStatus {
		if status.status == fileStatus {
			files = append(files, status.file)
		}
	}

	return files
}

func (d *Dropper) validateRootPath() error {
	rootAbs, err := filepath.Abs(d.rootPath.ToString())
	if err != nil {
		return fmt.Errorf("unable to get absolut path from %s", d.rootPath)
	}

	exists, isDir := d.fileFetcher.PathExists(domain.Path(rootAbs))
	if !exists {
		return fmt.Errorf("the path '%s' does not exists", d.rootPath)
	}

	if !isDir {
		return fmt.Errorf("the path priovided is not a directory")
	}

	d.rootPath = domain.Path(rootAbs)
	return nil
}
