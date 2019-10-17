package cli

import (
	"context"
	"fmt"
	"log"
	"time"
)

//go:generate moq -out dropper_mock_test.go . dropper
type dropper interface {
	SyncFiles(ctx context.Context) error
	SendUpdates(ctx context.Context)
}

type DropCLI struct {
	dropper dropper
}

func NewDropCLI(dropper dropper) *DropCLI {
	return &DropCLI{dropper: dropper}
}

func (d *DropCLI) Run(ctx context.Context) error {
	log.Println("Starting dropper...")
	if err := d.runSync(ctx); err != nil {
		return err
	}
	d.runUpdate(ctx)
	return nil
}

func (d *DropCLI) runSync(ctx context.Context) error {
	if err := d.dropper.SyncFiles(ctx); err != nil {
		return fmt.Errorf("unable to start sync cli: %w", err)
	}

	ticker := time.NewTicker(800 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := d.dropper.SyncFiles(ctx); err != nil {
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
