// Package pgx provides a PostgreSQL repository implementation backed by pgx.
package pgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neatflowcv/release-watcher/internal/pkg/domain"
	"github.com/neatflowcv/release-watcher/internal/pkg/repository"
)

const (
	createProjectsTableSQL = `
CREATE TABLE IF NOT EXISTS projects (
	name text NOT NULL,
	url text NOT NULL,
	detected_version text NOT NULL,
	acknowledged_version text NOT NULL
)`

	createProjectSQL = `
INSERT INTO projects (
	name,
	url,
	detected_version,
	acknowledged_version
) VALUES (
	$1,
	$2,
	$3,
	$4
)`

	listProjectsSQL = `
SELECT
	name,
	url,
	detected_version,
	acknowledged_version
FROM projects
ORDER BY name`
)

var _ repository.Repository = (*Repository)(nil)

// Repository stores projects in PostgreSQL through pgx.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a PostgreSQL repository from a PostgreSQL DSN.
func NewRepository(dsn string) (*Repository, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	_, err = pool.Exec(context.Background(), createProjectsTableSQL)
	if err != nil {
		pool.Close()

		return nil, fmt.Errorf("create projects table: %w", err)
	}

	return &Repository{pool: pool}, nil
}

// Close releases repository database connections.
func (r *Repository) Close() {
	r.pool.Close()
}

// CreateProject stores a project.
func (r *Repository) CreateProject(
	ctx context.Context,
	project *domain.Project,
) error {
	_, err := r.pool.Exec(
		ctx,
		createProjectSQL,
		project.Name(),
		project.URL(),
		project.DetectedVersion(),
		project.AcknowledgedVersion(),
	)
	if err != nil {
		return fmt.Errorf("create project: %w", err)
	}

	return nil
}

// ListProjects returns all projects ordered by name.
func (r *Repository) ListProjects(
	ctx context.Context,
) ([]*domain.Project, error) {
	rows, err := r.pool.Query(ctx, listProjectsSQL)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	projects := make([]*domain.Project, 0)

	for rows.Next() {
		var (
			name                string
			url                 string
			detectedVersion     string
			acknowledgedVersion string
		)

		scanErr := rows.Scan(
			&name,
			&url,
			&detectedVersion,
			&acknowledgedVersion,
		)
		if scanErr != nil {
			return nil, fmt.Errorf("scan project: %w", scanErr)
		}

		project := domain.NewProject(
			name,
			url,
			detectedVersion,
			acknowledgedVersion,
		)
		projects = append(projects, project)
	}

	rowsErr := rows.Err()
	if rowsErr != nil {
		return nil, fmt.Errorf("iterate projects: %w", rowsErr)
	}

	return projects, nil
}
