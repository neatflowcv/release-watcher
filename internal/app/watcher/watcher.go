// Package watcher periodically scans projects for release updates.
package watcher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/neatflowcv/release-watcher/internal/pkg/repository"
)

const defaultScanInterval = 6 * time.Hour

var (
	errMissingRepository = errors.New("missing project repository")
)

// Watcher periodically scans all projects.
type Watcher struct {
	repository repository.Repository
	scheduler  gocron.Scheduler
}

// NewWatcher creates a watcher for project release scans.
func NewWatcher(projectRepository repository.Repository) (*Watcher, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("create scheduler: %w", err)
	}

	watcher := &Watcher{
		repository: projectRepository,
		scheduler:  scheduler,
	}

	_, err = scheduler.NewJob(
		gocron.DurationJob(defaultScanInterval),
		gocron.NewTask(watcher.handle),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)
	if err != nil {
		return nil, fmt.Errorf("create scan job: %w", err)
	}

	return watcher, nil
}

// Run starts periodic full-project scans.
func (w *Watcher) Run() error {
	w.scheduler.Start()

	return nil
}

// Close stops the watcher.
func (w *Watcher) Close() {
	_ = w.scheduler.Shutdown()
}

func (w *Watcher) handle(ctx context.Context) error {
	_, err := w.repository.ListProjects(ctx)
	if err != nil {
		return fmt.Errorf("list projects: %w", err)
	}

	return nil
}
