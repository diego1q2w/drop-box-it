package cli

import (
	"context"
	"fmt"
	"github.com/diego1q2w/drop-box-it/pkg/drop/domain"
	"log"
	"time"
)

//go:generate moq -out dropper_mock_test.go . dropper
type dropper interface {
	SyncFiles(ctx context.Context, rootPath domain.Path) error
	SendUpdates(ctx context.Context)
}

type DropCLI struct {
	dropper dropper
}

func NewDropCLI(dropper dropper) *DropCLI {
	return &DropCLI{dropper: dropper}
}

func (d *DropCLI) Run(ctx context.Context, rootPath string) error {
	if err := d.runSync(ctx, rootPath); err != nil {
		return err
	}
	d.runUpdate(ctx)
	return nil
}

func (d *DropCLI) runSync(ctx context.Context, rootPath string) error {
	if err := d.dropper.SyncFiles(ctx, domain.Path(rootPath)); err != nil {
		return fmt.Errorf("unable to start sync cli: %w", err)
	}

	ticker := time.NewTicker(800 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := d.dropper.SyncFiles(ctx, domain.Path(rootPath)); err != nil {
					log.Printf("unexpected error while running sync cli: %s", err)
				}
			}
		}
	}()

	return nil
}

func (d *DropCLI) runUpdate(ctx context.Context) {
	ticker := time.NewTicker(2000 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				d.sendUpdate(ctx)
			}
		}
	}()
}

func (d *DropCLI) sendUpdate(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 1500*time.Millisecond)
	defer cancel()
	d.dropper.SendUpdates(ctx)
}
