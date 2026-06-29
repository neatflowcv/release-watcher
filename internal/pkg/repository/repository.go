// Package repository defines persistence boundaries for release watcher data.
package repository

import (
	"context"

	"github.com/neatflowcv/release-watcher/internal/pkg/domain"
)

// Repository stores and retrieves release watcher projects.
type Repository interface {
	CreateProject(ctx context.Context, project *domain.Project) error
	ListProjects(ctx context.Context) ([]*domain.Project, error)
}
