package app

import (
	"context"
	"fmt"
	"github.com/diego1q2w/drop-box-it/pkg/box/domain"
	"log"
	"path/filepath"
)

type hasher = func(s string) uint64

//go:generate moq -out file_boxer_mock_test.go . fileBoxer
type fileBoxer interface {
	WriteFile(file domain.File) error
	DeleteFile(path domain.Path) error
}

type Box struct {
	fileBoxer       fileBoxer
	hasher          hasher
	rootPath        string
	numberOfWorkers int
	channelWorkers  []chan fileWork
	ready           bool
}

type fileWork struct {
	action domain.Action
	path   domain.Path
	file   domain.File
}

func NewBox(fileBoxer fileBoxer, hasher hasher, rootPath string, workers int) *Box {
	b := &Box{fileBoxer: fileBoxer, rootPath: rootPath, numberOfWorkers: workers, ready: false, hasher: hasher}
	// This should run by an external manager that will enable graceful shutdown
	b.spinUpWorkers(context.Background())
	b.ready = true
	return b
}

func (b *Box) WriteDocuments(ctx context.Context, file domain.File) error {
	worker := b.hasher(file.Path.ToString()) % uint64(b.numberOfWorkers)
	return b.routeRequest(fileWork{
		action: domain.Write,
		file:   file,
	}, int(worker))
}

func (b *Box) DeleteDocuments(ctx context.Context, path domain.Path) error {
	worker := b.hasher(path.ToString()) % uint64(b.numberOfWorkers)
	return b.routeRequest(fileWork{
		action: domain.Delete,
		path:   path,
	}, int(worker))
}

func (b *Box) routeRequest(fileWork fileWork, worker int) error {
	if !b.ready {
		return fmt.Errorf("the boxer is not ready")
	}

	b.channelWorkers[worker] <- fileWork
	return nil
}

func (b *Box) spinUpWorkers(ctx context.Context) {
	var chans []chan fileWork
	for i := 0; i < b.numberOfWorkers; i++ {
		ch := b.spinWorker(ctx)
		chans = append(chans, ch)
	}

	b.channelWorkers = chans
}

func (b *Box) spinWorker(ctx context.Context) chan fileWork {
	ch := make(chan fileWork)

	go func() {
		for fWork := range ch {
			b.processDocument(fWork)
		}
	}()

	return ch
}

func (b *Box) processDocument(fWork fileWork) {
	switch fWork.action {
	case domain.Write:
		if err := b.writeDocument(fWork.file); err != nil {
			log.Printf("unable to write document: %s", err)
		}
	case domain.Delete:
		if err := b.deleteDocument(fWork.path); err != nil {
			log.Printf("unable to delete document: %s", err)
		}
	}
}

func (b *Box) writeDocument(file domain.File) error {
	log.Printf("writting document with path '%s'", file.Path)
	file.Path = domain.Path(filepath.Join(b.rootPath, file.Path.ToString()))
	return b.fileBoxer.WriteFile(file)
}

func (b *Box) deleteDocument(path domain.Path) error {
	log.Printf("deleting document with path '%s'", path)
	path = domain.Path(filepath.Join(b.rootPath, path.ToString()))
	return b.fileBoxer.DeleteFile(path)
}
